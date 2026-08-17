package convert_equal_to_deep_equal

import (
	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "Equal: use DeepEqual for slices" }

func Match() Matcher {
	return Matcher{
		"__a.Equal(__b, __array)": tapeImported,
		"__a.Equal(__array, __b)": tapeImported,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, __array)": "__a.DeepEqual(__b, __array)",
		"__a.Equal(__array, __b)": "__a.DeepEqual(__array, __b)",
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
