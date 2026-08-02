package convert_equal_to_deep_equal

import . "coderaiser/indra/types"

// Self is the plugin value used in engine-loader and Nested maps.
var Self = self{}

type self struct{}

func (self) Report() string    { return Report() }
func (self) Match() Matcher    { return Match() }
func (self) Replace() Replacer { return Replace() }

// Top-level exported funcs are readable and testable individually.

func Report() string { return "Equal: use DeepEqual for slices" }

func Match() Matcher {
	return Matcher{
		"__a.Equal(__b, __array)": nil,
		"__a.Equal(__array, __b)": nil,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b, __array)": "__a.DeepEqual(__b, __array)",
		"__a.Equal(__array, __b)": "__a.DeepEqual(__array, __b)",
	}
}
