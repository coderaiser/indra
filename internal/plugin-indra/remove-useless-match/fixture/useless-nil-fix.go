//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report() string { return "remove Test.Skip call" }

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
