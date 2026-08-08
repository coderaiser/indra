package plugin_indra

import (
	"coderaiser/indra/internal/plugin_indra/apply_fixture_name_to_message"
	"coderaiser/indra/internal/plugin_indra/convert_for_to_create"
	"coderaiser/indra/internal/plugin_indra/convert_inspect_to_traverse"
	"coderaiser/indra/internal/plugin_indra/convert_switch_to_if"
	"coderaiser/indra/internal/plugin_indra/remove_useless_match"
	"coderaiser/indra/internal/plugin_indra/replace_test_message"
	"coderaiser/indra/types"
)

// Rules returns the indra meta-rules. Enable via .indra.toml: "indra" = "on".
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "remove-useless-match", Plugin: remove_useless_match.Plugin{}},
		{Name: "convert-for-to-create-test", Plugin: convert_for_to_create.Plugin{}},
		{Name: "convert-switch-to-if", Plugin: convert_switch_to_if.Plugin{}},
		{Name: "apply-fixture-name-to-message", Plugin: apply_fixture_name_to_message.Plugin{}},
		{Name: "replace-test-message", Plugin: replace_test_message.Plugin{}},
		{Name: "convert-inspect-to-traverse", Plugin: convert_inspect_to_traverse.Plugin{}},
	}
}
