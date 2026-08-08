//go:build ignore

package fixture

import . "coderaiser/indra/types"

// non-composite-return: the body does not directly return a composite literal,
// so the Match is not the Matcher{...} shape and must be left alone.
func Match() Matcher {
	m := Matcher{}
	return m
}
