package convert_equal_to_not_ok

import (
	. "coderaiser/indra/types"

	"coderaiser/indra/internal/plugin_tape/tapeguard"
)

func Report() string { return "convert Equal(x, nil/false) to NotOk(x)" }

func Match() Matcher {
	return Matcher{
		"__a.Equal(__b, nil)":   tapeImported,
		"__a.Equal(__b, false)": tapeImported,
		"__a.Equal(__b, \"\")":  tapeImported,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, nil)":   "__a.NotOk(__b)",
		"__a.Equal(__b, false)": "__a.NotOk(__b)",
		"__a.Equal(__b, \"\")":  "__a.NotOk(__b)",
	}
}

// tapeImported is the per-rule guard: the pattern only fires inside a file
// that imports go-tape. Import detection delegates to tapeguard.
func tapeImported(vars Vars, path Path) bool { return tapeguard.Imported(vars, path) }

type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
