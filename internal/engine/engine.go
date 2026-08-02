// Package engine runs plugins against Go source.
//
// It is a shim over the higher-level pipeline: it parses the source, resolves
// a flat list of Self-shaped plugin values into concrete kinds, then runs
// pattern-based replacers and whole-AST traversers, applying Fix when asked.
//
// Internal engine depends on types; it does not import plugin packages.
package engine

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

	"coderaiser/indra/compare"
	"coderaiser/indra/types"
)

// Vars is an alias for the shared hole-bindings type.
type Vars = types.Vars

// Place is an alias for the shared lint-finding type.
type Place = types.Place

// MatchFn is an alias for the shared pattern guard type.
type MatchFn = types.MatchFn

// TraverseVisitor is an alias for the shared AST visitor type.
type TraverseVisitor = types.VisitFn

// replacerPlugin is implemented by any Self-shaped plugin exposing
// pattern-based rewriting (Report/Match/Replace). Only exported methods are
// used so reflection can detect the kind across packages.
type replacerPlugin interface {
	Report() string
	Match() types.Matcher
	Replace() types.Replacer
}

// traverserPlugin is implemented by any Self-shaped plugin exposing
// whole-node analysis (Report/Traverse/Fix). Only exported methods are used
// so reflection can detect the kind across packages.
type traverserPlugin interface {
	Report() string
	Traverse() types.Traverser
	Fix(node ast.Node, places []types.Place)
}

// resolvedPlugin is a Self-shaped plugin normalized into a concrete kind.
type resolvedPlugin struct {
	name      string
	kind      string // "replacer" or "traverser"
	report    func() string
	replacer  replacerPlugin
	traverser traverserPlugin
}

// load normalizes a flat list of Self-shaped values, expanding Nested groups.
// It panics on an unknowable plugin value.
func load(plugins []any) []resolvedPlugin {
	var out []resolvedPlugin
	for _, p := range plugins {
		if nested, ok := p.(types.Nested); ok {
			for name, sub := range nested {
				out = append(out, resolve(name, sub))
			}
			continue
		}
		out = append(out, resolve(deriveName(p), p))
	}
	return out
}

// resolve detects a plugin's kind by method presence.
func resolve(name string, p any) resolvedPlugin {
	switch v := p.(type) {
	case replacerPlugin:
		return resolvedPlugin{name: name, kind: "replacer", report: v.Report, replacer: v}
	case traverserPlugin:
		return resolvedPlugin{name: name, kind: "traverser", report: v.Report, traverser: v}
	default:
		panic("unknown plugin kind: " + name)
	}
}

// deriveName gets a display name from the type's name, else the package path.
func deriveName(p any) string {
	t := reflect.TypeOf(p)
	name := t.Name()
	if name == "" {
		name = t.String()
	}
	return name
}

// Indra runs plugins against src. It returns the (possibly modified) source,
// the collected places, and a parse error when the source is invalid.
func Indra(src []byte, plugins []any, fix bool) ([]byte, []types.Place, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return src, nil, err
	}

	var places []types.Place
	resolved := load(plugins)

	// Traverser plugins own whole-file analysis and may fix in place.
	for _, p := range resolved {
		if p.kind != "traverser" {
			continue
		}
		for key, visitor := range p.traverser.Traverse() {
			switch key {
			case "*ast.File":
				findings := visitor(file, Vars{})
				if fix && len(findings) > 0 {
					p.traverser.Fix(file, findings)
				}
				for _, pl := range findings {
					pl.Rule = p.name
					places = append(places, pl)
				}
			case "*ast.BlockStmt":
				ast.Inspect(file, func(n ast.Node) bool {
					block, ok := n.(*ast.BlockStmt)
					if !ok {
						return true
					}
					findings := visitor(block, Vars{})
					if fix && len(findings) > 0 {
						p.traverser.Fix(block, findings)
					}
					for _, pl := range findings {
						pl.Rule = p.name
						places = append(places, pl)
					}
					return true
				})
			}
		}
	}

	// Pattern-based (replacer) plugins run over every statement.
	var rewrites []rewrite
	walkStmts(file, func(stmt ast.Stmt, block *ast.BlockStmt, idx int) {
		for _, p := range resolved {
			if p.kind != "replacer" {
				continue
			}
			patterns := p.replacer.Match()
			for pattern, guard := range patterns {
				vars := compare.Compare(stmt, pattern)
				if vars == nil {
					continue
				}
				if guard != nil && !guard(vars) {
					continue
				}
				places = append(places, types.Place{
					Rule:    p.name,
					Message: p.report(),
					Pos:     fset.Position(stmt.Pos()),
				})
				if fix {
					if tmpl, ok := p.replacer.Replace()[pattern]; ok {
						newStmts := substituteAndParse(tmpl, vars)
						if newStmts == nil {
							continue
						}
						rewrites = append(rewrites, rewrite{
							block:    block,
							idx:      idx,
							newStmts: newStmts,
						})
					}
				}
			}
		}
	})

	if fix && len(rewrites) > 0 {
		applyRewrites(rewrites)
		var buf bytes.Buffer
		_ = format.Node(&buf, fset, file)
		return buf.Bytes(), places, nil
	}

	return src, places, nil
}

// rewrite is a pending statement replacement.
type rewrite struct {
	block    *ast.BlockStmt
	idx      int
	newStmts []ast.Stmt
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
