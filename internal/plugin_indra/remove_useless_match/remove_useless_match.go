package remove_useless_match

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ ast.Node) string { return "remove useless Match" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessMatch}
}

func findUselessMatch(node ast.Node, push func(ast.Node)) {
	file := node.(*ast.File)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Match" {
			continue
		}
		if isUselessMatch(fn) {
			push(file)
			break
		}
	}
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
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			return false
		}
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
	if fn.Recv != nil || fn.Type == nil || fn.Type.Results == nil {
		return nil
	}
	if len(fn.Type.Results.List) != 1 || len(fn.Type.Results.List[0].Names) != 0 {
		return nil
	}
	ident, ok := fn.Type.Results.List[0].Type.(*ast.Ident)
	if !ok || ident.Name != "Matcher" {
		return nil
	}
	if fn.Body == nil || len(fn.Body.List) != 1 {
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

// Fix removes useless Match() decls from file.Decls. node is *ast.File;
// options is unused.
func Fix(node ast.Node, _ map[string]any) {
	file := node.(*ast.File)
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
