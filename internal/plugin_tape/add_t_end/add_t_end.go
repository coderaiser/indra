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

// missingEnd is a guard that accepts a test body which does not already end
// with t.End(). The path argument is unused but kept to satisfy MatchFn.
func missingEnd(vars Vars, _ Path) bool {
	return !stmtsContainEnd(vars["__body"].(BodySlice).Stmts)
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
