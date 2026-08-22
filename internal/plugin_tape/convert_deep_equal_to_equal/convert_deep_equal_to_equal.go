package convert_deep_equal_to_equal

import (
	. "coderaiser/indra/operator"
	. "coderaiser/indra/types"
)

func Report() string { return "Use Equal() when comparing primitives" }

func Match() Matcher {
	return Matcher{
		"__a.DeepEqual(__b, __c)":      isPrimitiveLit,
		"__a.DeepEqual(__b, __c, __d)": isPrimitiveLit,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.DeepEqual(__b, __c)":      "__a.Equal(__b, __c)",
		"__a.DeepEqual(__b, __c, __d)": "__a.Equal(__b, __c, __d)",
	}
}

// isPrimitiveLit guards to primitive operands, for which Equal and DeepEqual
// are interchangeable.
func isPrimitiveLit(vars Vars, _ Path) bool {
	return IsPrimitive(vars["__c"])
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
