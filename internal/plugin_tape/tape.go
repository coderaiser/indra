package tape

import (
	"coderaiser/indra/internal/plugin_tape/add_t_end"
	"coderaiser/indra/internal/plugin_tape/convert_equal_to_deep_equal"
	"coderaiser/indra/internal/plugin_tape/convert_equal_to_not_ok"
	"coderaiser/indra/internal/plugin_tape/convert_equal_to_ok"
	"coderaiser/indra/internal/plugin_tape/convert_no_error_to_not_ok"
	"coderaiser/indra/internal/plugin_tape/convert_ok_to_not_ok"
	"coderaiser/indra/internal/plugin_tape/extract_result_from_assertion"
	"coderaiser/indra/internal/plugin_tape/remove_skip"
	"coderaiser/indra/internal/plugin_tape/remove_useless_condition"
	"coderaiser/indra/internal/plugin_tape/remove_useless_prefix"
	"coderaiser/indra/types"
)

// Rules returns the tape sub-rules. The engine-loader expands "tape" into
// "tape/remove-skip", "tape/add-t-end", etc. Each rule's shape is decided by
// its Plugin struct, not by this list.
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "remove-skip", Plugin: remove_skip.Plugin{}},
		{Name: "add-t-end", Plugin: add_t_end.Plugin{}},
		{Name: "convert-equal-to-deep-equal", Plugin: convert_equal_to_deep_equal.Plugin{}},
		{Name: "convert-equal-to-ok", Plugin: convert_equal_to_ok.Plugin{}},
		{Name: "convert-equal-to-not-ok", Plugin: convert_equal_to_not_ok.Plugin{}},
		{Name: "convert-ok-to-not-ok", Plugin: convert_ok_to_not_ok.Plugin{}},
		{Name: "extract-result-from-assertion", Plugin: extract_result_from_assertion.Plugin{}},
		{Name: "remove-useless-prefix", Plugin: remove_useless_prefix.Plugin{}},
		{Name: "remove-useless-condition", Plugin: remove_useless_condition.Plugin{}},
		{Name: "convert-no-error-to-not-ok", Plugin: convert_no_error_to_not_ok.Plugin{}},
	}
}
