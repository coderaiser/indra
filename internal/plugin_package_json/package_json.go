package plugin_package_json

import (
	"coderaiser/indra/internal/plugin_package_json/rename_version"
	"coderaiser/indra/types"
)

// Rules returns the package-json sub-rules. The engine-loader expands
// "package-json" into "package-json/rename-version", etc.
func Rules() []types.Rule {
	return []types.Rule{
		{Name: "rename-version", Plugin: rename_version.Plugin{}},
	}
}
