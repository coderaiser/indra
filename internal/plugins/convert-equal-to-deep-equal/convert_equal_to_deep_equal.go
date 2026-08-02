package convert_equal_to_deep_equal

import "coderaiser/indra/internal/engine"

var Plugin = engine.Plugin{
	Name:    "convert-equal-to-deep-equal",
	Report:  report,
	Match:   match,
	Replace: replace,
}

func report() string { return "Equal: use DeepEqual for slices" }

func match() map[string]engine.MatchFn {
	return map[string]engine.MatchFn{
		"__a.Equal(__b, __array)": nil,
		"__a.Equal(__array, __b)": nil,
	}
}

func replace() map[string]string {
	return map[string]string{
		"__a.Equal(__b, __array)": "__a.DeepEqual(__b, __array)",
		"__a.Equal(__array, __b)": "__a.DeepEqual(__array, __b)",
	}
}
