package plugins

import (
	. "coderaiser/indra/types"
	"coderaiser/indra/internal/plugins/add-t-end"
	"coderaiser/indra/internal/plugins/convert-equal-to-deep-equal"
	"coderaiser/indra/internal/plugins/extract-result-from-assertion"
	"coderaiser/indra/internal/plugins/remove-skip"
	"coderaiser/indra/internal/plugins/remove-unused-import"
	"coderaiser/indra/internal/plugins/remove-unused-variable"
	"coderaiser/indra/internal/plugins/tape"
)

// Plugins is the ordered list of top-level plugin package paths.
// Order matters: convert-equal-to-deep-equal before extract-result-from-assertion.
// engine-loader loads each via go/packages at init and detects kind by exported funcs.
var Plugins = []string{
	"coderaiser/indra/internal/plugins/remove-skip",
	"coderaiser/indra/internal/plugins/convert-equal-to-deep-equal",
	"coderaiser/indra/internal/plugins/add-t-end",
	"coderaiser/indra/internal/plugins/extract-result-from-assertion",
	"coderaiser/indra/internal/plugins/remove-unused-import",
	"coderaiser/indra/internal/plugins/remove-unused-variable",
}

// NestedPlugins maps group name → Nested rules map.
// engine-loader expands each into "group/rule" entries.
var NestedPlugins = map[string]Nested{
	"tape": tape.Rules,
}

// TODO: temporary shim deleted when internal/lint is rewritten in the
// engine-loader/engine-runner/engine-processor refactor.
var All = []any{
	remove_skip.Self,
	convert_equal_to_deep_equal.Self,
	add_t_end.Self,
	extract_result_from_assertion.Self,
	remove_unused_import.Self,
	remove_unused_variable.Self,
}
