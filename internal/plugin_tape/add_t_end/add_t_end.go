package add_t_end

import (
	"go/ast"

	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "tape: missing t.End()" }

func Match() Matcher {
	return Matcher{
		`Test(__a, __b, func(__a *Test.T) { __body })`:      tapeAndMissingEnd,
		`Test.Only(__a, __b, func(__a *Test.T) { __body })`: tapeAndMissingEnd,
	}
}

func Replace() Replacer {
	return Replacer{
		`Test(__a, __b, func(__a *Test.T) { __body })`:      "Test(__a, __b, func(__a *Test.T) {\n__body\n__a.End()\n})",
		`Test.Only(__a, __b, func(__a *Test.T) { __body })`: "Test.Only(__a, __b, func(__a *Test.T) {\n__body\n__a.End()\n})",
	}
}

// tapeImported is the per-rule guard: the pattern only fires inside a file
// that imports go-tape. Import detection delegates to tapeguard.
func tapeImported(vars Vars, path Path) bool { return tapeguard.Imported(vars, path) }

// tapeAndMissingEnd accepts only tape files whose matched body is missing an
// End() call.
func tapeAndMissingEnd(vars Vars, path Path) bool {
	return tapeImported(vars, path) && missingEnd(vars, path)
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
// with t.End().
func missingEnd(vars Vars, _ Path) bool {
	body, ok := vars["__body"].(BodySlice)
	if !ok {
		return false
	}
	return !stmtsContainEnd(body.Stmts)
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
