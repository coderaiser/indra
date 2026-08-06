package main

import (
	engine_loader         "coderaiser/indra/engine_loader"
	plugin_indra          "coderaiser/indra/internal/plugin_indra"
	plugin_tape           "coderaiser/indra/internal/plugin_tape"
	remove_unused_import  "coderaiser/indra/internal/remove_unused_import"
	remove_unused_variable "coderaiser/indra/internal/remove_unused_variable"
)

// Registry is the single source of truth for all top-level plugins.
// Sub-plugins are owned by their group's Rules() func — not listed here.
// Adding a new standalone plugin: one line here.
// Adding a new sub-plugin: one entry in its group's Rules() + Plugin{} struct.
var Registry = []engine_loader.PluginFuncs{
	{Name: "tape", Rules: plugin_tape.Rules()},
	{Name: "indra", Rules: plugin_indra.Rules()},
	{Name: "remove-unused-import", Plugin: remove_unused_import.Plugin{}},
	{Name: "remove-unused-variable", Plugin: remove_unused_variable.Plugin{}},
}
