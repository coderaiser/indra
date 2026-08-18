package add_t_end

import (
	"go/ast"

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

// missingEnd is a guard that accepts a test body which does not already end
// with an End() call. The [match] config already scopes tape rules to
// *_test.go files, so no import guard is needed here.
func missingEnd(vars Vars, _ Path) bool {
	body, ok := vars["__body"].(BodySlice)
	if !ok {
		return false
	}
	return !stmtsContainEnd(body.Stmts)
}

// stmtsContainEnd reports whether any statement in stmts is a call to an End
// method (t.End()).
func stmtsContainEnd(stmts []ast.Stmt) bool {
	for _, s := range stmts {
		if Compare(s, "__.End()") {
			return true
		}
	}
	return false
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
