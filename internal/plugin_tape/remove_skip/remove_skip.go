package remove_skip

import (
	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "remove Test.Skip call" }

func Match() Matcher {
	return Matcher{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: tapeImported,
	}
}

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}

// tapeImported is the per-rule guard: the pattern only fires inside a file
// that imports go-tape. Import detection delegates to tapeguard.
func tapeImported(vars Vars, path Path) bool { return tapeguard.Imported(vars, path) }

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
