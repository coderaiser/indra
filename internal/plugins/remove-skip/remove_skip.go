package remove_skip

import . "coderaiser/indra/types"

// Self is the plugin value used in engine-loader and Nested maps.
var Self = self{}

type self struct{}

func (self) Report() string    { return Report() }
func (self) Match() Matcher    { return Match() }
func (self) Replace() Replacer { return Replace() }

// Top-level exported funcs are readable and testable individually.

func Report() string { return "remove Test.Skip call" }

func Match() Matcher {
	return Matcher{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: nil,
	}
}

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}
