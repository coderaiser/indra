package convert_ok_to_not_ok

import (
	. "coderaiser/indra/types"
)

// Top-level exported funcs are readable and testable individually.

func Report() string { return "convert Ok to NotOk" }

func Replace() Replacer {
	return Replacer{
		"__a.Ok(__b == nil)": "__a.NotOk(__b)",
		"__a.Ok(!__b)":       "__a.NotOk(__b)",
	}
}

// Plugin wraps the rule for the registry: a replacer. The [match] config
// already scopes tape rules to *_test.go files, so no per-plugin import guard
// is needed.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
