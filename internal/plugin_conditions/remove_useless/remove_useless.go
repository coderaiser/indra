package remove_useless

import . "coderaiser/indra/types"

func Report() string { return "Avoid useless condition" }

func Replace() Replacer {
	return Replacer{
		"__a != nil": "__b",
		"__a == nil": "!__a",
	}
}

// Plugin wraps the rule for the registry: a replacer without a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
