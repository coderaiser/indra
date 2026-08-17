package convert_equal_to_ok

import (
	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

func Report() string { return "convert Equal(x, true) to Ok(x)" }

func Match() Matcher {
	return Matcher{
		"__a.Equal(__b, true)": tapeImported,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, true)": "__a.Ok(__b)",
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
