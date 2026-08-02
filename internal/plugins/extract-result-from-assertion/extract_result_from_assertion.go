package extractresultfromassertion

import "coderaiser/indra/internal/engine"

var Plugin = engine.Plugin{
	Name:    "extract-result-from-assertion",
	Report:  report,
	Match:   match,
	Replace: replace,
}

func report() string { return "extract inline expression from assertion" }

func match() map[string]engine.MatchFn {
	return map[string]engine.MatchFn{
		"__a.Equal(__b(__args), __c)":     nil,
		"__a.DeepEqual(__b(__args), __c)": nil,
		"__a.Equal(__b, __array)":         nil,
		"__a.DeepEqual(__b, __array)":     nil,
	}
}

func replace() map[string]string {
	return map[string]string{
		"__a.Equal(__b(__args), __c)":     "result := __b(__args)\n__a.Equal(result, __c)",
		"__a.DeepEqual(__b(__args), __c)": "result := __b(__args)\n__a.DeepEqual(result, __c)",
		"__a.Equal(__b, __array)":         "expected := __array\n__a.Equal(__b, expected)",
		"__a.DeepEqual(__b, __array)":     "expected := __array\n__a.DeepEqual(__b, expected)",
	}
}
