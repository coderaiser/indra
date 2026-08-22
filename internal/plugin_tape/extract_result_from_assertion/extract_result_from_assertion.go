package extract_result_from_assertion

import (
	. "coderaiser/indra/operator"
	. "coderaiser/indra/types"
)

func Report() string { return "extract inline expression from assertion" }

// Match guards the call-extraction patterns so a function-call result is not
// re-extracted when a "result" variable is already declared in scope (which
// would shadow the injected declaration), and the composite-literal patterns
// so "expected" is not re-declared when it already exists.
func Match() Matcher {
	return Matcher{
		"__a.Equal(__b(__args), __c)":     noResultInScope,
		"__a.DeepEqual(__b(__args), __c)": noResultInScope,

		"__a.Equal(__b, __array)":      expectedNotDeclared,
		"__a.DeepEqual(__b, __array)":  expectedNotDeclared,
		"__a.Equal(__b, __struct)":     expectedNotDeclared,
		"__a.DeepEqual(__b, __struct)": expectedNotDeclared,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b(__args), __c)":     "result := __b(__args)\n__a.Equal(result, __c)",
		"__a.DeepEqual(__b(__args), __c)": "result := __b(__args)\n__a.DeepEqual(result, __c)",
		"__a.Equal(__b, __array)":         "expected := __array\n__a.Equal(__b, expected)",
		"__a.DeepEqual(__b, __array)":     "expected := __array\n__a.DeepEqual(__b, expected)",
		"__a.Equal(__b, __struct)":        "expected := __struct\n__a.Equal(__b, expected)",
		"__a.DeepEqual(__b, __struct)":    "expected := __struct\n__a.DeepEqual(__b, expected)",
	}
}

// noResultInScope is a guard that rejects re-extraction when a "result"
// variable is already declared in scope (which would shadow the injected
// declaration). Statement matching always runs inside a block.
func noResultInScope(_ Vars, path Path) bool {
	return GetBinding(path, "result") == nil
}

// expectedNotDeclared rejects composite-literal extraction when an "expected"
// variable is already declared in scope.
func expectedNotDeclared(_ Vars, path Path) bool {
	return GetBinding(path, "expected") == nil
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
