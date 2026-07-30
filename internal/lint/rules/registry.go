package rules

import "coderaiser/indra/internal/lint/rule"

var All = []rule.Rule{
	&AssertCount{},
	&NoSkip{},
	&NoEqualSlice{},
	&RequireTEnd{},
}
