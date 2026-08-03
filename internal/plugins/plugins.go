package plugins

import (
	engine_loader "coderaiser/indra/engine-loader"
	plugin_indra "coderaiser/indra/internal/plugin-indra"
	remove_useless_match "coderaiser/indra/internal/plugin-indra/remove-useless-match"
	"coderaiser/indra/internal/plugins/add-t-end"
	"coderaiser/indra/internal/plugins/convert-equal-to-deep-equal"
	"coderaiser/indra/internal/plugins/convert-equal-to-not-ok"
	"coderaiser/indra/internal/plugins/extract-result-from-assertion"
	"coderaiser/indra/internal/plugins/remove-skip"
	"coderaiser/indra/internal/plugins/remove-unused-import"
	"coderaiser/indra/internal/plugins/remove-unused-variable"
	"coderaiser/indra/internal/plugins/tape"
)

// All is the single, ordered plugin registry. The tape group is the canonical
// owner of its sub-rules; those sub-plugins are not registered standalone here.
// Order matters: tape group runs convert-equal-to-deep-equal before
// extract-result-from-assertion (both match Equal+array patterns).
var All = []engine_loader.PluginFuncs{
	{Name: "tape", Path: "coderaiser/indra/internal/plugins/tape", Rules: tape.Rules},
	{Name: "remove-unused-import", Path: "coderaiser/indra/internal/plugins/remove-unused-import", Report: remove_unused_import.Report, Traverse: remove_unused_import.Traverse, Fix: remove_unused_import.Fix},
	{Name: "remove-unused-variable", Path: "coderaiser/indra/internal/plugins/remove-unused-variable", Report: remove_unused_variable.Report, Traverse: remove_unused_variable.Traverse, Fix: remove_unused_variable.Fix},
	{Name: "remove-useless-match", Path: "coderaiser/indra/internal/plugin-indra/remove-useless-match", Report: remove_useless_match.Report, Replace: remove_useless_match.Replace},
	{Name: "indra", Path: "coderaiser/indra/internal/plugin-indra", Rules: plugin_indra.Rules},
}

// Providers holds the PluginFuncs for tape sub-rules so loader.Load can expand
// the tape group into "tape/*" rules, without registering them standalone.
var Providers = []engine_loader.PluginFuncs{
	{Name: "remove-skip", Path: "coderaiser/indra/internal/plugins/remove-skip", Report: remove_skip.Report, Replace: remove_skip.Replace},
	{Name: "convert-equal-to-deep-equal", Path: "coderaiser/indra/internal/plugins/convert-equal-to-deep-equal", Report: convert_equal_to_deep_equal.Report, Match: convert_equal_to_deep_equal.Match, Replace: convert_equal_to_deep_equal.Replace},
	{Name: "convert-equal-to-not-ok", Path: "coderaiser/indra/internal/plugins/convert-equal-to-not-ok", Report: convert_equal_to_not_ok.Report, Match: convert_equal_to_not_ok.Match, Replace: convert_equal_to_not_ok.Replace},
	{Name: "add-t-end", Path: "coderaiser/indra/internal/plugins/add-t-end", Report: add_t_end.Report, Match: add_t_end.Match, Replace: add_t_end.Replace},
	{Name: "extract-result-from-assertion", Path: "coderaiser/indra/internal/plugins/extract-result-from-assertion", Report: extract_result_from_assertion.Report, Match: extract_result_from_assertion.Match, Replace: extract_result_from_assertion.Replace},
}

// LoadInput is the slice loader.Load needs to resolve top-level rules and the
// nested tape group. Providers are appended so their paths resolve inside the
// group, though they never surface as standalone rules.
func LoadInput() []engine_loader.PluginFuncs {
	return append(append([]engine_loader.PluginFuncs{}, All...), Providers...)
}


