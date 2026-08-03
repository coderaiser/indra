package plugins

import (
	engine_loader "coderaiser/indra/engine-loader"
	plugin_indra "coderaiser/indra/internal/plugin-indra"
	remove_useless_match "coderaiser/indra/internal/plugin-indra/remove-useless-match"
	"coderaiser/indra/internal/plugins/add-t-end"
	"coderaiser/indra/internal/plugins/convert-equal-to-deep-equal"
	"coderaiser/indra/internal/plugins/extract-result-from-assertion"
	"coderaiser/indra/internal/plugins/remove-skip"
	"coderaiser/indra/internal/plugins/remove-unused-import"
	"coderaiser/indra/internal/plugins/remove-unused-variable"
	"coderaiser/indra/internal/plugins/tape"
)

// All is the single, ordered plugin registry.
// Order matters: convert-equal-to-deep-equal before extract-result-from-assertion.
// Nested plugins (tape) carry their Rules map; the loader expands them into
// "group/rule" entries at load time.
var All = []engine_loader.PluginFuncs{
	{Name: "remove-skip", Path: "coderaiser/indra/internal/plugins/remove-skip", Report: remove_skip.Report, Match: remove_skip.Match, Replace: remove_skip.Replace},
	{Name: "convert-equal-to-deep-equal", Path: "coderaiser/indra/internal/plugins/convert-equal-to-deep-equal", Report: convert_equal_to_deep_equal.Report, Match: convert_equal_to_deep_equal.Match, Replace: convert_equal_to_deep_equal.Replace},
	{Name: "add-t-end", Path: "coderaiser/indra/internal/plugins/add-t-end", Report: add_t_end.Report, Match: add_t_end.Match, Replace: add_t_end.Replace},
	{Name: "extract-result-from-assertion", Path: "coderaiser/indra/internal/plugins/extract-result-from-assertion", Report: extract_result_from_assertion.Report, Match: extract_result_from_assertion.Match, Replace: extract_result_from_assertion.Replace},
	{Name: "remove-unused-import", Path: "coderaiser/indra/internal/plugins/remove-unused-import", Report: remove_unused_import.Report, Traverse: remove_unused_import.Traverse, Fix: remove_unused_import.Fix},
	{Name: "remove-unused-variable", Path: "coderaiser/indra/internal/plugins/remove-unused-variable", Report: remove_unused_variable.Report, Traverse: remove_unused_variable.Traverse, Fix: remove_unused_variable.Fix},
	{Name: "tape", Path: "coderaiser/indra/internal/plugins/tape", Rules: tape.Rules},
	{Name: "remove-useless-match", Path: "coderaiser/indra/internal/plugin-indra/remove-useless-match", Report: remove_useless_match.Report, Match: remove_useless_match.Match, Replace: remove_useless_match.Replace},
	{Name: "indra", Path: "coderaiser/indra/internal/plugin-indra", Rules: plugin_indra.Rules},
}

