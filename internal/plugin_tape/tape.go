package tape

import . "coderaiser/indra/types"

// Rules groups tape-related sub-plugins.
// engine-loader expands "tape" → "tape/remove-skip", "tape/add-t-end", etc.
var Rules = Nested{
	"remove-skip":                   "coderaiser/indra/internal/plugin_tape/remove_skip",
	"add-t-end":                     "coderaiser/indra/internal/plugin_tape/add_t_end",
	"convert-equal-to-deep-equal":   "coderaiser/indra/internal/plugin_tape/convert_equal_to_deep_equal",
	"convert-equal-to-not-ok":       "coderaiser/indra/internal/plugin_tape/convert_equal_to_not_ok",
	"convert-ok-to-not-ok":          "coderaiser/indra/internal/plugin_tape/convert_ok_to_not_ok",
	"extract-result-from-assertion": "coderaiser/indra/internal/plugin_tape/extract_result_from_assertion",
	"remove-useless-condition":      "coderaiser/indra/internal/plugin_tape/remove_useless_condition",
	"remove-useless-prefix":         "coderaiser/indra/internal/plugin_tape/remove_useless_prefix",
	"convert-no-error-to-not-ok":    "coderaiser/indra/internal/plugin_tape/convert_no_error_to_not_ok",
}
