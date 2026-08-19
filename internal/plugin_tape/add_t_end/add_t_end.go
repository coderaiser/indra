package add_t_end

import (
	"go/ast"

	. "coderaiser/indra/operator"
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

// missingEnd is a guard that accepts a test body which does not already contain
// an End() call (whether as a bare statement, an assigned form, or implied by a
// trailing callback argument). The [match] config already scopes tape rules to
// *_test.go files, so no import guard is needed here.
func missingEnd(vars Vars, _ Path) bool {
	body := vars["__body"].(BodySlice)
	stmts := body.Stmts
	if len(stmts) == 0 {
		return true
	}
	// t.End() anywhere in body (putout: compareAny)
	if CompareAny("__a.End()", stmts) {
		return false
	}
	// const result = t.End() assigned form
	if CompareAny("__a := __b.End()", stmts) {
		return false
	}
	last := stmts[len(stmts)-1]
	// last stmt is a call whose last arg is a func literal → callback pattern
	if expr, ok := last.(*ast.ExprStmt); ok {
		if call, ok := expr.X.(*ast.CallExpr); ok && len(call.Args) > 0 {
			if IsFuncLit(call.Args[len(call.Args)-1]) {
				return false
			}
		}
	}
	return true
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
