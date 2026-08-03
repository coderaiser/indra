package extract_result_from_assertion

import . "coderaiser/indra/types"

// Top-level exported funcs are readable and testable individually.

func Report() string { return "extract inline expression from assertion" }

func Match() Matcher {
	return Matcher{
		"__a.Equal(__b(__args), __c)":     nil,
		"__a.DeepEqual(__b(__args), __c)": nil,
		"__a.Equal(__b, __array)":         nil,
		"__a.DeepEqual(__b, __array)":     nil,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b(__args), __c)":     "result := __b(__args)\n__a.Equal(result, __c)",
		"__a.DeepEqual(__b(__args), __c)": "result := __b(__args)\n__a.DeepEqual(result, __c)",
		"__a.Equal(__b, __array)":         "expected := __array\n__a.Equal(__b, expected)",
		"__a.DeepEqual(__b, __array)":     "expected := __array\n__a.DeepEqual(__b, expected)",
	}
}
