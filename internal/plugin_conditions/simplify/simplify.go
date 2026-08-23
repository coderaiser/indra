// Package simplify replaces an if/else whose branches are identical with the
// branch body itself: the condition is irrelevant when both sides do the same
// thing. The repeated __b hole binds both branches, so the engine's pattern
// matcher enforces the equality. The optional-chaining patterns of putout's
// simplify are JS-specific and do not apply to Go.
package simplify

import (
	. "coderaiser/indra/types"
)

func Report() string { return "simplify condition" }

func Replace() Replacer {
	return Replacer{
		"if __a { __b } else { __b }": "__b",
	}
}

// Plugin wraps the rule for the registry: a replacer.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
