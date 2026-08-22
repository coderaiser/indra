package remove_default_messages

import (
	. "coderaiser/indra/types"
)

func Report() string { return "Avoid passing default messages to operators" }

func Replace() Replacer {
	return Replacer{
		`__a.Ok(__b, "should be truthy")`:                     "__a.Ok(__b)",
		`__a.NotOk(__b, "should be falsy")`:                   "__a.NotOk(__b)",
		`__a.Match(__b, __c, "should match")`:                 "__a.Match(__b, __c)",
		`__a.NotMatch(__b, __c, "should not match")`:          "__a.NotMatch(__b, __c)",
		`__a.Equal(__b, __c, "should equal")`:                 "__a.Equal(__b, __c)",
		`__a.NotEqual(__b, __c, "should not equal")`:          "__a.NotEqual(__b, __c)",
		`__a.DeepEqual(__b, __c, "should deep equal")`:        "__a.DeepEqual(__b, __c)",
		`__a.NotDeepEqual(__b, __c, "should not deep equal")`: "__a.NotDeepEqual(__b, __c)",
	}
}

// Plugin wraps the rule for the registry: a pure replacer without a guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Replace() Replacer { return Replace() }
