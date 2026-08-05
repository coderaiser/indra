package remove_useless_match

import "go/ast"

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
