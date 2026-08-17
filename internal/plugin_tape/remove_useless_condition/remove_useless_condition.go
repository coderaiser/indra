package remove_useless_condition

import (
	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

// Package remove_useless_condition simplifies tape assertions whose condition
// is redundant: Ok(x != nil) means the same as Ok(x), and the easier form
// reads better.

func Report() string { return "remove useless condition" }

func Match() Matcher {
	return Matcher{
		"__a.Ok(__b != nil)":      tapeImported,
		"__a.NotOk(__b == nil)":   tapeImported,
		"__a.Ok(__b != false)":    tapeImported,
		"__a.NotOk(__b == false)": tapeImported,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Ok(__b != nil)":      "__a.Ok(__b)",
		"__a.NotOk(__b == nil)":   "__a.NotOk(__b)",
		"__a.Ok(__b != false)":    "__a.Ok(__b)",
		"__a.NotOk(__b == false)": "__a.NotOk(__b)",
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
