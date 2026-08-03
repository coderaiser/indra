//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report() string { return "some plugin" }

func Replace() Replacer {
	return Replacer{`some pattern`: "replacement"}
}
