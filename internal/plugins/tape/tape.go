package tape

import . "coderaiser/indra/types"

// Rules groups tape-related sub-plugins.
// engine-loader expands "tape" → "tape/remove-skip", "tape/add-t-end", etc.
var Rules = Nested{
	"remove-skip":                   "coderaiser/indra/internal/plugins/remove_skip",
	"add-t-end":                     "coderaiser/indra/internal/plugins/add_t_end",
	"convert-equal-to-deep-equal":   "coderaiser/indra/internal/plugins/convert_equal_to_deep_equal",
	"convert-equal-to-not-ok":       "coderaiser/indra/internal/plugins/convert_equal_to_not_ok",
	"convert-ok-to-not-ok":          "coderaiser/indra/internal/plugins/convert_ok_to_not_ok",
	"extract-result-from-assertion": "coderaiser/indra/internal/plugins/extract_result_from_assertion",
	"remove-useless-condition":      "coderaiser/indra/internal/plugins/remove_useless_condition",
	"remove-useless-prefix":         "coderaiser/indra/internal/plugins/remove_useless_prefix",
	"convert-no-error-to-not-ok":    "coderaiser/indra/internal/plugins/convert_no_error_to_not_ok",
}
