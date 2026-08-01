package addtend

import (
	"go/ast"

	"coderaiser/indra/internal/engine"
	"coderaiser/indra/internal/engine/compare"
)

var Plugin = engine.Plugin{
	Name:   "add-t-end",
	Report: func() string { return "tape: missing t.End()" },
	Match: func() map[string]engine.MatchFn {
		guard := func(vars engine.Vars) bool {
			body, ok := vars["__body"].(compare.BodySlice)
			if !ok {
				return false
			}
			return !stmtsContainEnd(body.Stmts)
		}
		return map[string]engine.MatchFn{
			`tape.Test(__t, __name, func(__t *tape.T) { __body })`: guard,
			`tape.Only(__t, __name, func(__t *tape.T) { __body })`: guard,
		}
	},
	Replace: func() map[string]string {
		return map[string]string{
			`tape.Test(__t, __name, func(__t *tape.T) { __body })`: "tape.Test(__t, __name, func(__t *tape.T) {\n__body\n__t.End()\n})",
			`tape.Only(__t, __name, func(__t *tape.T) { __body })`: "tape.Only(__t, __name, func(__t *tape.T) {\n__body\n__t.End()\n})",
		}
	},
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
