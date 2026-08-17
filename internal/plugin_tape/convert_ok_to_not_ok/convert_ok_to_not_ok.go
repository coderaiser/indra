package convert_ok_to_not_ok

import (
	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "convert Ok to NotOk" }

func Match() Matcher {
	return Matcher{
		"__a.Ok(__b == nil)": tapeImported,
		"__a.Ok(!__b)":       tapeImported,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Ok(__b == nil)": "__a.NotOk(__b)",
		"__a.Ok(!__b)":       "__a.NotOk(__b)",
	}
}

// tapeImported is the per-rule guard: the pattern only fires inside a file
// that imports go-tape. Import detection delegates to tapeguard so it lives in
// one place and stays at full coverage.
func tapeImported(vars Vars, path Path) bool { return tapeguard.Imported(vars, path) }

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
