//go:build ignore

package fixture

import . "coderaiser/indra/types"

func Report(_ Path) string { return "partial" }

func Replace() Replacer {
	return Replacer{}
}
