package plugin_indra

import (
	"coderaiser/indra/internal/plugin_indra/convert_for_to_create_test"
	"coderaiser/indra/internal/plugin_indra/remove_useless_match"
	"coderaiser/indra/types"
)

// Rules returns the indra meta-rules. All rules are Disabled by default;
// enable them in .indra.toml: "indra" = "on".
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "remove-useless-match", Plugin: remove_useless_match.Plugin{}, Disabled: true},
		{Name: "convert-for-to-create-test", Plugin: convert_for_to_create_test.Plugin{}, Disabled: true},
	}
}
