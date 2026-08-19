// Package convert_no_error_to_not_ok rewrites tape's t.NoError(err) into the
// equivalent t.NotOk(err). Only checks with the tape receiver name t are
// considered, so the collision-prone NoError name is never rewritten elsewhere.
package convert_no_error_to_not_ok

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report() string { return "convert NoError(err) to NotOk(err)" }

// Match guards on the receiver name: only a receiver identifier literally
// named t (the go-tape convention) is rewritten. Mirrors putout's receiver
// check in convert-equal-to-ok et al.
func Match() Matcher {
	return Matcher{
		"__a.NoError(__b)": func(vars Vars, _ Path) bool {
			ident, ok := vars["__a"].(*ast.Ident)
			return ok && ident.Name == "t"
		},
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.NoError(__b)": "__a.NotOk(__b)",
	}
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }

