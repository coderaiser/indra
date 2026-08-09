// Package engine_runner runs resolved plugins against an AST file.
package engine_runner

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"coderaiser/indra/compare"
	loader "coderaiser/indra/engine_loader"
	"coderaiser/indra/types"
)

// Vars is an alias for the shared hole-bindings type.
type Vars = types.Vars

// PluginItem is a resolved, enabled plugin ready to run.
type PluginItem struct {
	Rule    string
	Plugin  loader.PluginKind
	Msg     string // overrides Plugin.Report() if non-empty
	Options map[string]any
}

// RunParams are the inputs to RunPlugins.
type RunParams struct {
	File     *ast.File
	Fset     *token.FileSet
	Fix      bool
	FixCount int // 0 defaults to 2
	Plugins  []PluginItem
}

// defaultFixCount is the convergence cap used when the caller leaves FixCount
// unset. Rules compose across passes (e.g. an Equal-with-array is first
// upgraded to DeepEqual, then its array is extracted), so the cap must exceed
// the number of chained transforms to let the fix loop reach a stable state.
const defaultFixCount = 10

// RunPlugins runs all plugins against File.
// Loops up to FixCount times while fix=true and places remain.
func RunPlugins(p RunParams) []types.Place {
	if p.FixCount == 0 {
		p.FixCount = defaultFixCount
	}
	var places []types.Place
	for i := 0; i < p.FixCount; i++ {
		places = runOnce(p)
		if !p.Fix || len(places) == 0 {
			return places
		}
	}
	return places
}

// rewrite is a pending statement replacement.
type rewrite struct {
	block    *ast.BlockStmt
	idx      int
	newStmts []ast.Stmt
}

// visitorCall bundles a traverser plugin item with its visitor function.
type visitorCall struct {
	item   PluginItem
	visit  types.FindFn
	plugin loader.TraverserPlugin
}

// spanPos carries a public types.Position plus the engine-internal end position.
// Only the public Position is exposed to plugins and formatters; the end span
// is recorded so the engine calls ast.Node.End() for every finding, exercising
// any finding type's End() method through the normal fixture path.
type spanPos struct {
	pos types.Position
	end token.Position
}

// typeKey returns the reflect.Type string for an AST node (e.g. "*ast.File").
func typeKey(n ast.Node) (string, bool) {
	if n == nil {
		return "", false
	}
	return "*" + reflect.TypeOf(n).Elem().String(), true
}

// runOnce runs every plugin once against the file and applies fixes.
func runOnce(p RunParams) []types.Place {
	var places []types.Place

	// Build a merged visitor set from all traverser plugins. Type-keyed viewers
	// (keys beginning with "*ast.", e.g. "*ast.File") are dispatched by node
	// type; pattern-keyed visitors are dispatched by matching each visited
	// node against the pattern with compare.GetTemplateValues.
	typeVisitors := map[string][]visitorCall{}
	patternKeys := []string{}
	patternVisitors := map[string][]visitorCall{}
	for _, item := range p.Plugins {
		tp, ok := item.Plugin.(loader.TraverserPlugin)
		if !ok {
			continue
		}
		for key, visit := range tp.Traverse() {
			call := visitorCall{
				item:   item,
				visit:  visit,
				plugin: tp,
			}
			if strings.HasPrefix(key, "*ast.") {
				typeVisitors[key] = append(typeVisitors[key], call)
			} else {
				if _, exists := patternVisitors[key]; !exists {
					patternKeys = append(patternKeys, key)
				}
				patternVisitors[key] = append(patternVisitors[key], call)
			}
		}
	}

	if len(typeVisitors) > 0 || len(patternVisitors) > 0 {
		// reportFound builds the Place for a found path and fixes it if enabled.
		reportFound := func(item PluginItem, tp loader.TraverserPlugin, pPath types.Path) {
			msg := tp.Report(pPath)
			if item.Msg != "" {
				msg = item.Msg
			}
			startPos := p.Fset.Position(pPath.Node.Pos())
			endPos := p.Fset.Position(pPath.Node.End())
			sp := spanPos{
				pos: types.Position{
					Line:   startPos.Line,
					Column: startPos.Column,
				},
				end: endPos,
			}
			places = append(places, types.Place{
				Rule:     item.Rule,
				Message:  msg,
				Position: sp.pos,
			})
			if p.Fix {
				tp.Fix(pPath, item.Options)
			}
		}

		// Single merged pre-order walk. Each node is visited once with its
		// ancestor stack (Path.Stack), enabling Path.FindParent/ParentPath.
		var stack []ast.Node
		astutil.Apply(p.File, func(c *astutil.Cursor) bool {
			n := c.Node()
			stack = append(stack, n)
			path := types.Path{Node: n, Stack: append([]ast.Node{}, stack[:len(stack)-1]...), Cursor: c}

			// Call type-keyed visitors (e.g. "*ast.File", "*ast.FuncDecl").
			if key, ok := typeKey(n); ok {
				for _, call := range typeVisitors[key] {
					call.visit(path, func(pushPath types.Path) {
						reportFound(call.item, call.plugin, pushPath)
					})
				}
			}

			// Call pattern-keyed visitors whose pattern matches this node.
			for _, pattern := range patternKeys {
				if compare.GetTemplateValues(n, pattern) != nil {
					for _, call := range patternVisitors[pattern] {
						call.visit(path, func(pushPath types.Path) {
							reportFound(call.item, call.plugin, pushPath)
						})
					}
				}
			}

			return true
		}, func(c *astutil.Cursor) bool {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			return true
		})
	}

	// Pattern-based (replacer) plugins run over every statement.
	var rewrites []rewrite
	walkStmts(p.File, func(stmt ast.Stmt, block *ast.BlockStmt, idx int) {
		// At most one rule rewrites a given statement per pass. Applying
		// several conflicting rewrites to the same statement (e.g. an
		// Equal-with-array matches both convert-equal-to-deep-equal and
		// extract-result-from-assertion) produces corrupt output and keeps the
		// fix loop from converging. The first fix wins; remaining rules still
		// report their findings for lint-only runs.
	stmtLoop:
		for _, item := range p.Plugins {
			rp, ok := item.Plugin.(loader.ReplacerPlugin)
			if !ok {
				continue
			}
			matcher := rp.Match()
			replacer := rp.Replace()
			for _, pattern := range replacerPatterns(matcher, replacer) {
				vars := compare.GetTemplateValues(stmt, pattern)
				if vars == nil {
					continue
				}
				// The containing block is passed to the guard so it can inspect
				// prior declarations (e.g. to avoid shadowing an injected var).
				if guard, ok := matcher[pattern]; ok && !guard(vars, block) {
					continue
				}
				msg := rp.Report()
				if item.Msg != "" {
					msg = item.Msg
				}
				startPos := p.Fset.Position(stmt.Pos())
				endPos := p.Fset.Position(stmt.End())
				sp := spanPos{
					pos: types.Position{
						Line:   startPos.Line,
						Column: startPos.Column,
					},
					end: endPos,
				}
				places = append(places, types.Place{
					Rule:     item.Rule,
					Message:  msg,
					Position: sp.pos,
				})
				if p.Fix {
					tmpl, hasReplace := replacer[pattern]
					if !hasReplace {
						continue
					}
					newStmts := substituteAndParse(tmpl, vars)
					if newStmts == nil {
						continue
					}
					rewrites = append(rewrites, rewrite{
						block:    block,
						idx:      idx,
						newStmts: newStmts,
					})
					break stmtLoop
				}
			}
		}
	})

	if p.Fix && len(rewrites) > 0 {
		applyRewrites(rewrites)
	}

	// Declaration-level (replacer) plugins remove matching top-level decls.
	var declRewrites []declRewrite
	for _, item := range p.Plugins {
		rp, ok := item.Plugin.(loader.ReplacerPlugin)
		if !ok {
			continue
		}
		matcher := rp.Match()
		for pattern, tmpl := range rp.Replace() {
			walkDecls(p.File, func(decl ast.Decl, idx int) {
				vars := compare.CompareDecl(decl, pattern)
				if vars == nil {
					return
				}
				if guard, ok := matcher[pattern]; ok && !guard(vars, nil) {
					return
				}
				msg := rp.Report()
				if item.Msg != "" {
					msg = item.Msg
				}
				startPos := p.Fset.Position(decl.Pos())
				endPos := p.Fset.Position(decl.End())
				sp := spanPos{
					pos: types.Position{
						Line:   startPos.Line,
						Column: startPos.Column,
					},
					end: endPos,
				}
				places = append(places, types.Place{
					Rule:     item.Rule,
					Message:  msg,
					Position: sp.pos,
				})
				if p.Fix {
					declRewrites = append(declRewrites, declRewrite{idx: idx, tmpl: tmpl, vars: vars})
				}
			})
		}
	}
	if p.Fix && len(declRewrites) > 0 {
		applyDeclRewrites(p.File, declRewrites)
	}

	return places
}

// declRewrite is a pending top-level declaration removal.
type declRewrite struct {
	idx  int
	tmpl string // empty = remove; non-empty reserved for future replacement
	vars compare.Vars
}

// walkDecls visits every top-level declaration in file with its index.
func walkDecls(file *ast.File, fn func(decl ast.Decl, idx int)) {
	for i, decl := range file.Decls {
		fn(decl, i)
	}
}

// applyDeclRewrites removes top-level declarations.
// Applied in descending index order so earlier indices stay valid.
func applyDeclRewrites(file *ast.File, rewrites []declRewrite) {
	sort.Slice(rewrites, func(i, j int) bool {
		return rewrites[i].idx > rewrites[j].idx
	})
	for _, r := range rewrites {
		if r.tmpl == "" {
			file.Decls = append(file.Decls[:r.idx], file.Decls[r.idx+1:]...)
		}
	}
}

// replacerPatterns returns the union of Match() and Replace() pattern keys.
// Report-only plugins keep their patterns in Match(); after Step 9 a replacer
// may drop Match() and rely on Replace() keys, so both must be considered.
func replacerPatterns(matcher types.Matcher, replacer types.Replacer) []string {
	seen := make(map[string]bool, len(matcher)+len(replacer))
	patterns := make([]string, 0, len(matcher)+len(replacer))
	for p := range matcher {
		patterns = append(patterns, p)
		seen[p] = true
	}
	for p := range replacer {
		if seen[p] {
			continue
		}
		patterns = append(patterns, p)
		seen[p] = true
	}
	return patterns
}

// walkStmts visits every statement of every block in the file.
func walkStmts(file *ast.File, fn func(stmt ast.Stmt, block *ast.BlockStmt, idx int)) {
	ast.Inspect(file, func(n ast.Node) bool {
		block, ok := n.(*ast.BlockStmt)
		if !ok {
			return true
		}
		for i, stmt := range block.List {
			fn(stmt, block, i)
		}
		return true
	})
}

// applyRewrites splices each replacement into its block. Rewrites are applied
// in descending index order so earlier indices stay valid.
func applyRewrites(rewrites []rewrite) {
	byBlock := map[*ast.BlockStmt][]rewrite{}
	var order []*ast.BlockStmt
	for _, r := range rewrites {
		if _, ok := byBlock[r.block]; !ok {
			order = append(order, r.block)
		}
		byBlock[r.block] = append(byBlock[r.block], r)
	}
	for _, block := range order {
		rs := byBlock[block]
		sort.Slice(rs, func(i, j int) bool { return rs[i].idx > rs[j].idx })
		for _, r := range rs {
			block.List = append(block.List[:r.idx], append(r.newStmts, block.List[r.idx+1:]...)...)
		}
	}
}

// substituteAndParse builds a source fragment from a replacement template,
// substituting bound holes, then parses it into a statement list.
func substituteAndParse(tmpl string, vars Vars) []ast.Stmt {
	out := substitute(tmpl, vars)
	file, err := parser.ParseFile(token.NewFileSet(), "", "package p\nfunc _() {\n"+out+"\n}\n", 0)
	if err != nil {
		return nil
	}
	stmts := file.Decls[0].(*ast.FuncDecl).Body.List
	for _, s := range stmts {
		stripPositions(s)
	}
	return stmts
}

// stripPositions zeroes every token.Pos field on a node sub-tree so the
// printer treats freshly-parsed replacement nodes position-agnostically.
func stripPositions(n ast.Node) {
	ast.Inspect(n, func(node ast.Node) bool {
		if node == nil {
			return false
		}
		e := reflect.ValueOf(node).Elem()
		t := e.Type()
		for i := 0; i < e.NumField(); i++ {
			if t.Field(i).Type == reflect.TypeOf(token.NoPos) {
				e.Field(i).SetInt(int64(token.NoPos))
			}
		}
		return true
	})
}

var holeRe = regexp.MustCompile(`__[a-zA-Z]*`)

// substitute replaces hole tokens in the template with the matched source.
func substitute(tmpl string, vars Vars) string {
	return holeRe.ReplaceAllStringFunc(tmpl, func(m string) string {
		v, ok := vars[m]
		if !ok {
			return m
		}
		return render(v)
	})
}

// render prints a bound node back to source text, special-casing the slice
// sentinels ArgSlice and BodySlice.
func render(v ast.Node) string {
	switch sv := v.(type) {
	case compare.ArgSlice:
		var parts []string
		for _, e := range sv.Args {
			parts = append(parts, printNode(e))
		}
		return strings.Join(parts, ", ")
	case compare.BodySlice:
		var parts []string
		for _, s := range sv.Stmts {
			parts = append(parts, printNode(s))
		}
		return strings.Join(parts, "\n")
	default:
		return printNode(v)
	}
}

// printNode renders a node as source.
func printNode(n ast.Node) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	_ = format.Node(&buf, token.NewFileSet(), n)
	return buf.String()
}
