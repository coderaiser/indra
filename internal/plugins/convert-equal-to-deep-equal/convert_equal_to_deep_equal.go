package convertequaltoDeepEqual

import "coderaiser/indra/internal/engine"

var Plugin = engine.Plugin{
	Name:   "convert-equal-to-deep-equal",
	Report: func() string { return "Equal: use DeepEqual for slices" },
	Match: func() map[string]engine.MatchFn {
		return map[string]engine.MatchFn{
			"__recv.Equal(__a, __array)": nil,
			"__recv.Equal(__array, __b)": nil,
		}
	},
	Replace: func() map[string]string {
		return map[string]string{
			"__recv.Equal(__a, __array)": "__recv.DeepEqual(__a, __array)",
			"__recv.Equal(__array, __b)": "__recv.DeepEqual(__array, __b)",
		}
	},
}
