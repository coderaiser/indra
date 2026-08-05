package plugins

import (
	engine_loader "coderaiser/indra/engine_loader"
	plugin_indra "coderaiser/indra/internal/plugin_indra"
	remove_useless_match "coderaiser/indra/internal/plugin_indra/remove_useless_match"
	"coderaiser/indra/internal/plugins/add_t_end"
	"coderaiser/indra/internal/plugins/convert_equal_to_deep_equal"
	"coderaiser/indra/internal/plugins/convert_equal_to_not_ok"
	convert_no_error_to_not_ok "coderaiser/indra/internal/plugins/convert_no_error_to_not_ok"
	"coderaiser/indra/internal/plugins/convert_ok_to_not_ok"
	"coderaiser/indra/internal/plugins/extract_result_from_assertion"
	"coderaiser/indra/internal/plugins/remove_skip"
	"coderaiser/indra/internal/plugins/remove_unused_import"
	"coderaiser/indra/internal/plugins/remove_unused_variable"
	"coderaiser/indra/internal/plugins/remove_useless_condition"
	"coderaiser/indra/internal/plugins/remove_useless_prefix"
	"coderaiser/indra/internal/plugins/tape"
)

// All is the single, ordered plugin registry. The tape group is the canonical
// owner of its sub-rules; those sub-plugins are not registered standalone here.
// Order matters: tape group runs convert-equal-to-deep-equal before
// extract-result-from-assertion (both match Equal+array patterns).
var All = []engine_loader.PluginFuncs{
	{Name: "tape", Path: "coderaiser/indra/internal/plugins/tape", Rules: tape.Rules},
	{Name: "remove-unused-import", Path: "coderaiser/indra/internal/plugins/remove_unused_import", Report: remove_unused_import.Report, Traverse: remove_unused_import.Traverse, Fix: remove_unused_import.Fix},
	{Name: "convert-no-error-to-not-ok", Path: "coderaiser/indra/internal/plugins/convert_no_error_to_not_ok", Report: convert_no_error_to_not_ok.Report, Traverse: convert_no_error_to_not_ok.Traverse, Fix: convert_no_error_to_not_ok.Fix},
	{Name: "remove-unused-variable", Path: "coderaiser/indra/internal/plugins/remove_unused_variable", Report: remove_unused_variable.Report, Traverse: remove_unused_variable.Traverse, Fix: remove_unused_variable.Fix},
	{Name: "remove-useless-match", Path: "coderaiser/indra/internal/plugin_indra/remove_useless_match", Report: remove_useless_match.Report, Traverse: remove_useless_match.Traverse, Fix: remove_useless_match.Fix},
	{Name: "indra", Path: "coderaiser/indra/internal/plugin_indra", Rules: plugin_indra.Rules},
}

// Providers holds the PluginFuncs for tape sub-rules so loader.Load can expand
// the tape group into "tape/*" rules, without registering them standalone.
var Providers = []engine_loader.PluginFuncs{
	{Name: "remove-skip", Path: "coderaiser/indra/internal/plugins/remove_skip", Report: remove_skip.Report, Replace: remove_skip.Replace},
	{Name: "convert-equal-to-deep-equal", Path: "coderaiser/indra/internal/plugins/convert_equal_to_deep_equal", Report: convert_equal_to_deep_equal.Report, Replace: convert_equal_to_deep_equal.Replace},
	{Name: "convert-equal-to-not-ok", Path: "coderaiser/indra/internal/plugins/convert_equal_to_not_ok", Report: convert_equal_to_not_ok.Report, Replace: convert_equal_to_not_ok.Replace},
	{Name: "convert-ok-to-not-ok", Path: "coderaiser/indra/internal/plugins/convert_ok_to_not_ok", Report: convert_ok_to_not_ok.Report, Replace: convert_ok_to_not_ok.Replace},
	{Name: "add-t-end", Path: "coderaiser/indra/internal/plugins/add_t_end", Report: add_t_end.Report, Match: add_t_end.Match, Replace: add_t_end.Replace},
	{Name: "extract-result-from-assertion", Path: "coderaiser/indra/internal/plugins/extract_result_from_assertion", Report: extract_result_from_assertion.Report, Match: extract_result_from_assertion.Match, Replace: extract_result_from_assertion.Replace},
	{Name: "remove-useless-condition", Path: "coderaiser/indra/internal/plugins/remove_useless_condition", Report: remove_useless_condition.Report, Replace: remove_useless_condition.Replace},
	{Name: "remove-useless-prefix", Path: "coderaiser/indra/internal/plugins/remove_useless_prefix", Report: remove_useless_prefix.Report, Traverse: remove_useless_prefix.Traverse, Fix: remove_useless_prefix.Fix},
	{Name: "convert-no-error-to-not-ok", Path: "coderaiser/indra/internal/plugins/convert_no_error_to_not_ok", Report: convert_no_error_to_not_ok.Report, Traverse: convert_no_error_to_not_ok.Traverse, Fix: convert_no_error_to_not_ok.Fix},
}

// LoadInput is the slice loader.Load needs to resolve top-level rules and the
// nested tape group. Providers are appended so their paths resolve inside the
// group, though they never surface as standalone rules.
func LoadInput() []engine_loader.PluginFuncs {
	return append(append([]engine_loader.PluginFuncs{}, All...), Providers...)
}
