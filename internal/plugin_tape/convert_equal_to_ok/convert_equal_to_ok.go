package convert_equal_to_ok

import (
	. "coderaiser/indra/types"
)

func Report() string { return "convert Equal(x, true) to Ok(x)" }

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, true)":          "__a.Ok(__b)",
		"__a.Equal(__b, true, __c)":     "__a.Ok(__b, __c)",
		"__a.DeepEqual(__b, true)":      "__a.Ok(__b)",
		"__a.DeepEqual(__b, true, __c)": "__a.Ok(__b, __c)",
	}
}

type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
