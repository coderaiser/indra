// Package conditions groups rule plugins that apply to any assertion
// library — nothing here is tape-specific. Rules in this group run on all
// files by default, no glob restriction.
package conditions

import (
	"coderaiser/indra/internal/plugin_conditions/remove_useless_condition"
	"coderaiser/indra/types"
)

// Rules returns the conditions sub-rules. The engine-loader expands
// "conditions" into "conditions/remove-useless-condition", etc.
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "remove-useless-condition", Plugin: remove_useless_condition.Plugin{}},
	}
}
