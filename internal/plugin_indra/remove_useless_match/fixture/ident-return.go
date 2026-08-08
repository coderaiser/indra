//go:build ignore

package fixture

import . "coderaiser/indra/types"

// ident-return: the return value is an identifier, not a Matcher literal, so
// this is not the plugin Match shape and must be left alone.
func Match() Matcher {
	return someMatcher
}
