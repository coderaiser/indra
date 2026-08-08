//go:build ignore

package fixture

import . "coderaiser/indra/types"

// non-return-body: the single body statement is not a return, so this is not
// the plugin Match shape and must be left alone.
func Match() Matcher {
	_ = Matcher{}
}
