package switch_expected_with_result

import (
	. "coderaiser/indra/types"
)

func Report() string { return `"result" should be before "expected"` }

func Replace() Replacer {
	return Replacer{
		"__a.Equal(expected, __b)":          "__a.Equal(__b, expected)",
		"__a.Equal(expected, __b, __c)":     "__a.Equal(__b, expected, __c)",
		"__a.DeepEqual(expected, __b)":      "__a.DeepEqual(__b, expected)",
		"__a.DeepEqual(expected, __b, __c)": "__a.DeepEqual(__b, expected, __c)",
	}
}

// Plugin wraps the rule for the registry: a pure replacer without a guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
