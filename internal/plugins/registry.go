package plugins

import (
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
var All = []any{
	remove_skip.Self,
	convert_equal_to_deep_equal.Self,
	add_t_end.Self,
	extract_result_from_assertion.Self,
	remove_unused_import.Self,
	remove_unused_variable.Self,
}
