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

	"coderaiser/indra/internal/engine/compare"
)

// Vars is an alias for the compare hole-bindings map.
type Vars = compare.Vars

// MatchFn is an optional guard run after a pattern matches. It may inspect
// the bound holes and decide whether the match should be accepted.
type MatchFn = func(Vars) bool

// TraverseVisitor runs against a whole-file node (e.g. *ast.File) and
// returns any places it detected.
type TraverseVisitor = func(node ast.Node, vars Vars) []Place

// Plugin describes a single lint rule.
type Plugin struct {
	Name     string
	Report   func() string
	Match    func() map[string]MatchFn // nil = always match
	Replace  func() map[string]string  // nil = report only
	Traverse func() map[string]TraverseVisitor
}

// Place is a single lint finding.
type Place struct {
	Rule    string
	Message string
	Pos     token.Position
}

// Indra runs plugins against src. It returns the (possibly modified) source,
// the collected places, and a parse error when the source is invalid.
func Indra(src []byte, plugins []Plugin, fix bool) ([]byte, []Place, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return src, nil, err
	}

	var places []Place

	// Traverse plugins own whole-file analysis.
	for _, p := range plugins {
		if p.Traverse == nil {
			continue
		}
		for key, visitor := range p.Traverse() {
			switch key {
			case "*ast.File":
				for _, pl := range visitor(file, Vars{}) {
					pl.Rule = p.Name
					places = append(places, pl)
				}
			case "*ast.BlockStmt":
				ast.Inspect(file, func(n ast.Node) bool {
					block, ok := n.(*ast.BlockStmt)
					if !ok {
						return true
					}
					for _, pl := range visitor(block, Vars{}) {
						pl.Rule = p.Name
						places = append(places, pl)
					}
					return true
				})
			}
		}
	}

	// Pattern-based (reporter/replacer) plugins run over every statement.
	var rewrites []rewrite
	walkStmts(file, func(stmt ast.Stmt, block *ast.BlockStmt, idx int) {
		for _, p := range plugins {
			if p.Traverse != nil || p.Match == nil {
				continue
			}
			patterns := p.Match()
			for pattern, guard := range patterns {
				vars := compare.Compare(stmt, pattern)
				if vars == nil {
					continue
				}
				if guard != nil && !guard(vars) {
					continue
				}
				places = append(places, Place{
					Rule:    p.Name,
					Message: p.Report(),
					Pos:     fset.Position(stmt.Pos()),
				})
				if fix && p.Replace != nil {
					if tmpl, ok := p.Replace()[pattern]; ok {
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
