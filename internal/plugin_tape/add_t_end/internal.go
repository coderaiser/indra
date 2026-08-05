package add_t_end

import (
	"go/ast"

	"coderaiser/indra/compare"
	. "coderaiser/indra/types"
)

func stmtsContainEnd(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		expr, ok := s.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := expr.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if sel.Sel.Name == "End" {
			return true
		}
	}
	return false
}

// missingEnd is a guard that accepts a test body which does not already end
// with t.End(). The block argument is unused but kept to satisfy MatchFn.
func missingEnd(vars Vars, _ *ast.BlockStmt) bool {
	body, ok := vars["__body"].(compare.BodySlice)
	if !ok {
		return false
	}
	return !stmtsContainEnd(body.Stmts)
}
