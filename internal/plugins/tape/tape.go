package tape

import . "coderaiser/indra/types"

// Rules groups tape-related sub-plugins.
// engine-loader expands "tape" → "tape/remove-skip", "tape/add-t-end", etc.
var Rules = Nested{
	"remove-skip":                   "coderaiser/indra/internal/plugins/remove-skip",
	"add-t-end":                     "coderaiser/indra/internal/plugins/add-t-end",
	"convert-equal-to-deep-equal":   "coderaiser/indra/internal/plugins/convert-equal-to-deep-equal",
	"extract-result-from-assertion": "coderaiser/indra/internal/plugins/extract-result-from-assertion",
}
