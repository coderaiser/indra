package convert_ok_to_not_ok

import . "coderaiser/indra/types"

func Report() string { return "convert Ok(err == nil) to NotOk(err)" }

func Replace() Replacer {
	return Replacer{
		"__a.Ok(__b == nil)": "__a.NotOk(__b)",
	}
}

// Plugin wraps the rule for the registry: a replacer without a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
