// Package conditions groups rule plugins that apply to any assertion
// library — nothing here is tape-specific. Rules in this group run on all
// files by default, no glob restriction.
package conditions

import (
	"coderaiser/indra/internal/plugin_conditions/convert_switch_to_if"
	"coderaiser/indra/internal/plugin_conditions/remove_useless_comments"
	"coderaiser/indra/types"
)

// Rules returns the conditions sub-rules. The engine-loader expands
// "conditions" into "conditions/remove-useless-comments", etc.
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "remove-useless-comments", Plugin: remove_useless_comments.Plugin{}},
		{Name: "convert-switch-to-if", Plugin: convert_switch_to_if.Plugin{}},
	}
}
