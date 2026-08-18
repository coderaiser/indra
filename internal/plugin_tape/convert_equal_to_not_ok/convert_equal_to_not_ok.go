package convert_equal_to_not_ok

import (
	. "coderaiser/indra/types"
)

func Report() string { return "convert Equal(x, nil/false) to NotOk(x)" }

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, nil)":   "__a.NotOk(__b)",
		"__a.Equal(__b, false)": "__a.NotOk(__b)",
		"__a.Equal(__b, \"\")":  "__a.NotOk(__b)",
	}
}

type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
