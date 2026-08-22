package remove_useless_match

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove useless Match" }

// Fix removes useless Match() decls from file.Decls. node is *ast.File;
// options is unused.
func Fix(p Path, _ map[string]any) {
	file := p.Node.(*ast.File)
	kept := file.Decls[:0]
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name != nil && fn.Name.Name == "Match" && isUselessMatch(fn) {
			continue
		}
		kept = append(kept, decl)
	}
	file.Decls = kept
}

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessMatch}
}

// findUselessMatch pushes the file once when it declares a top-level Match()
// function that returns an empty Matcher or an all-nil-guard Matcher.
func findUselessMatch(p Path, push func(Path)) {
	p.Traverse(map[string]func(Path){
		"*ast.FuncDecl": func(declPath Path) {
			fn := declPath.Node.(*ast.FuncDecl)
			if fn.Name.Name != "Match" {
				return
			}
			if isUselessMatch(fn) {
				push(p)
				declPath.Stop()
			}
		},
	})
}

// isUselessMatch reports whether fn returns an empty Matcher or a Matcher
// whose every guard is nil. Either form makes Match() deletable: an empty
// Matcher does nothing, and all-nil guards are no-ops.
func isUselessMatch(fn *ast.FuncDecl) bool {
	lit := matcherLit(fn)
	if lit == nil {
		return false
	}
	if len(lit.Elts) == 0 {
		return true
	}
	for _, elt := range lit.Elts {
		kv := elt.(*ast.KeyValueExpr)
		ident, ok := kv.Value.(*ast.Ident)
		if !ok || ident.Name != "nil" {
			return false
		}
	}
	return true
}

// matcherLit returns the Matcher{...} composite literal returned by fn, or nil
// if fn does not have the func Match() Matcher { return Matcher{...} } shape.
func matcherLit(fn *ast.FuncDecl) *ast.CompositeLit {
	results := fn.Type.Results
	if results == nil || len(results.List) != 1 {
		return nil
	}
	ident, ok := results.List[0].Type.(*ast.Ident)
	if !ok || ident.Name != "Matcher" {
		return nil
	}
	if len(fn.Body.List) != 1 {
		return nil
	}
	retStmt, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(retStmt.Results) != 1 {
		return nil
	}
	lit, ok := retStmt.Results[0].(*ast.CompositeLit)
	if !ok {
		return nil
	}
	return lit
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
