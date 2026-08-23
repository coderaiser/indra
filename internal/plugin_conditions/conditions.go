// Package conditions groups rule plugins that apply to any assertion
// library — nothing here is tape-specific. Rules in this group run on all
// files by default, no glob restriction.
package conditions

import (
	"coderaiser/indra/internal/plugin_conditions/apply_early_return"
	"coderaiser/indra/internal/plugin_conditions/convert_switch_to_if"
	"coderaiser/indra/internal/plugin_conditions/merge_if_statements"
	"coderaiser/indra/internal/plugin_conditions/merge_if_with_else"
	"coderaiser/indra/internal/plugin_conditions/remove_boolean"
	"coderaiser/indra/internal/plugin_conditions/remove_useless_comments"
	"coderaiser/indra/internal/plugin_conditions/remove_useless_else"
	"coderaiser/indra/internal/plugin_conditions/reverse_condition"
	"coderaiser/indra/internal/plugin_conditions/simplify"
	"coderaiser/indra/types"
)

// Rules returns the conditions sub-rules. The engine-loader expands
// "conditions" into "conditions/remove-useless-comments", etc.
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "remove-useless-comments", Plugin: remove_useless_comments.Plugin{}},
		{Name: "convert-switch-to-if", Plugin: convert_switch_to_if.Plugin{}},
		{Name: "remove-boolean", Plugin: remove_boolean.Plugin{}},
		{Name: "reverse-condition", Plugin: reverse_condition.Plugin{}},
		{Name: "remove-useless-else", Plugin: remove_useless_else.Plugin{}},
		{Name: "merge-if-statements", Plugin: merge_if_statements.Plugin{}},
		{Name: "merge-if-with-else", Plugin: merge_if_with_else.Plugin{}},
		{Name: "apply-early-return", Plugin: apply_early_return.Plugin{}},
		{Name: "simplify", Plugin: simplify.Plugin{}},
	}
}
