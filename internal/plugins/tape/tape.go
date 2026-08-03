package tape

import . "coderaiser/indra/types"

// Rules groups tape-related sub-plugins.
// engine-loader expands "tape" → "tape/remove-skip", "tape/add-t-end".
var Rules = Nested{
	"remove-skip": "coderaiser/indra/internal/plugins/remove-skip",
	"add-t-end":   "coderaiser/indra/internal/plugins/add-t-end",
}
