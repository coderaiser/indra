package convert_equal_to_deep_equal

import (
	. "coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "Equal: use DeepEqual for slices" }

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, __array)": "__a.DeepEqual(__b, __array)",
		"__a.Equal(__array, __b)": "__a.DeepEqual(__array, __b)",
	}
}

// Plugin wraps the rule for the registry: a replacer. The [match] config
// already scopes tape rules to *_test.go files, so no per-plugin import guard
// is needed.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
