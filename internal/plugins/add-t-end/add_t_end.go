package add_t_end

import (
	"go/ast"

	"coderaiser/indra/compare"
	. "coderaiser/indra/types"
)

// Self is the plugin value used in engine-loader and Nested maps.
var Self = self{}

type self struct{}

func (self) Report() string    { return Report() }
func (self) Match() Matcher    { return Match() }
func (self) Replace() Replacer { return Replace() }

// Top-level exported funcs are readable and testable individually.

func Report() string { return "tape: missing t.End()" }

func Match() Matcher {
	guard := func(vars Vars) bool {
		body, ok := vars["__body"].(compare.BodySlice)
		if !ok {
			return false
		}
		return !stmtsContainEnd(body.Stmts)
	}
	return Matcher{
		`Test(__a, __b, func(__a *Test.T) { __body })`:      guard,
		`Test.Only(__a, __b, func(__a *Test.T) { __body })`: guard,
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
