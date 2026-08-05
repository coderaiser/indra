//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report() string { return "some plugin" }

func Match() Matcher {
	return Matcher{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: func(vars Vars) bool {
			return vars["__a"] != nil
		},
	}
}

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
