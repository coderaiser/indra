//go:build ignore

package fixture

import . "coderaiser/indra/types"

// remove-useless-match is the canonical happy path: a Match guard whose only
// entry is a nil guard (all-nil guards are no-ops, so the Match is deletable).
func Match() Matcher {
	return Matcher{
		`Test(__a, __b, func(__a *Test.T) { __body })`: nil,
	}
}
