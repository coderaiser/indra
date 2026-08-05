package add_t_end

import (
	"go/ast"

	"coderaiser/indra/compare"
	. "coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "tape: missing t.End()" }

func Match() Matcher {
	return Matcher{
		`Test(__a, __b, func(__a *Test.T) { __body })`:      missingEnd,
		`Test.Only(__a, __b, func(__a *Test.T) { __body })`: missingEnd,
	}
}

func Replace() Replacer {
	return Replacer{
		`Test(__a, __b, func(__a *Test.T) { __body })`:      "Test(__a, __b, func(__a *Test.T) {\n__body\n__a.End()\n})",
		`Test.Only(__a, __b, func(__a *Test.T) { __body })`: "Test.Only(__a, __b, func(__a *Test.T) {\n__body\n__a.End()\n})",
	}
}

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
