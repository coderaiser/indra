package extractresultfromassertion

import "coderaiser/indra/internal/engine"

var Plugin = engine.Plugin{
	Name:   "extract-result-from-assertion",
	Report: func() string { return "extract inline expression from assertion" },
	Match: func() map[string]engine.MatchFn {
		return map[string]engine.MatchFn{
			"__recv.Equal(__call(__args), __b)":     nil,
			"__recv.DeepEqual(__call(__args), __b)": nil,
			"__recv.Equal(__a, __array)":            nil,
			"__recv.DeepEqual(__a, __array)":        nil,
		}
	},
	Replace: func() map[string]string {
		return map[string]string{
			"__recv.Equal(__call(__args), __b)":     "result := __call(__args)\n__recv.Equal(result, __b)",
			"__recv.DeepEqual(__call(__args), __b)": "result := __call(__args)\n__recv.DeepEqual(result, __b)",
			"__recv.Equal(__a, __array)":            "expected := __array\n__recv.Equal(__a, expected)",
			"__recv.DeepEqual(__a, __array)":        "expected := __array\n__recv.DeepEqual(__a, expected)",
		}
	},
}
