package addtend

import (
	"go/ast"

	"coderaiser/indra/internal/engine"
	"coderaiser/indra/internal/engine/compare"
)

var Plugin = engine.Plugin{
	Name:    "add-t-end",
	Report:  report,
	Match:   match,
	Replace: replace,
}

func report() string { return "tape: missing t.End()" }

func match() map[string]engine.MatchFn {
	guard := func(vars engine.Vars) bool {
		body, ok := vars["__body"].(compare.BodySlice)
		if !ok {
			return false
		}
		return !stmtsContainEnd(body.Stmts)
	}
	return map[string]engine.MatchFn{
		`Test(__a, __b, func(__a *Test.T) { __body })`:      guard,
		`Test.Only(__a, __b, func(__a *Test.T) { __body })`: guard,
	}
}

func replace() map[string]string {
	return map[string]string{
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
