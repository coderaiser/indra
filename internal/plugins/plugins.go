package plugins

import (
	engine_loader "coderaiser/indra/engine_loader"
	plugin_indra "coderaiser/indra/internal/plugin_indra"
	plugin_tape "coderaiser/indra/internal/plugin_tape"
	remove_unused_import "coderaiser/indra/internal/remove_unused_import"
	remove_unused_variable "coderaiser/indra/internal/remove_unused_variable"
)

// Registry is the single, ordered plugin registry. No entry describes a
// plugin's shape (Report/Match/Replace/Traverse/Fix): a group carries Rules(),
// a leaf carries a Plugin{} struct, and the engine-loader resolves each kind
// through reflection. A plugin can switch between replacer and traverser
// without touching this list.
var Registry = []engine_loader.PluginFuncs{
	{Name: "tape", Rules: plugin_tape.Rules()},
	{Name: "indra", Rules: plugin_indra.Rules()},
	{Name: "remove-unused-import", Plugin: remove_unused_import.Plugin{}},
	{Name: "remove-unused-variable", Plugin: remove_unused_variable.Plugin{}},
}
