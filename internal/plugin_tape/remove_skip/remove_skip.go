package remove_skip

import (
	. "coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "remove Test.Skip call" }

func Replace() Replacer {
	return Replacer{
		`Test.Skip(__a, __b, func(__a *Test.T) { __body })`: "Test(__a, __b, func(__a *Test.T) {\n__body\n})",
	}
}

// Plugin wraps the rule for the registry: a replacer. The [match] config
// already scopes tape rules to *_test.go files, so no per-plugin import guard
// is needed.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
