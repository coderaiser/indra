package plugins

import (
	"coderaiser/indra/internal/engine"
	"coderaiser/indra/internal/plugins/add-t-end"
	"coderaiser/indra/internal/plugins/convert-equal-to-deep-equal"
	"coderaiser/indra/internal/plugins/extract-result-from-assertion"
	"coderaiser/indra/internal/plugins/remove-skip"
	"coderaiser/indra/internal/plugins/remove-unused-import"
	"coderaiser/indra/internal/plugins/remove-unused-variable"
)

// All is the ordered list of plugins run by indra.
// Order matters: convert-equal-to-deep-equal must run before
// extract-result-from-assertion (both match Equal+array patterns).
var All = []engine.Plugin{
	remove_skip.Plugin,
	convert_equal_to_deep_equal.Plugin,
	add_t_end.Plugin,
	extract_result_from_assertion.Plugin,
	remove_unused_import.Plugin,
	remove_unused_variable.Plugin,
}
