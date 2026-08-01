package plugins

import (
	"coderaiser/indra/internal/engine"
	addtend "coderaiser/indra/internal/plugins/add-t-end"
	convertequaltoDeepEqual "coderaiser/indra/internal/plugins/convert-equal-to-deep-equal"
	extractresult "coderaiser/indra/internal/plugins/extract-result-from-assertion"
	removeunusedimport "coderaiser/indra/internal/plugins/remove-unused-import"
	removeunusedvariable "coderaiser/indra/internal/plugins/remove-unused-variable"
	removeskip "coderaiser/indra/internal/plugins/remove-skip"
)

// All is the ordered list of plugins run by indra.
// Order matters: convert-equal-to-deep-equal must run before
// extract-result-from-assertion (both match Equal+array patterns).
var All = []engine.Plugin{
	removeskip.Plugin,
	convertequaltoDeepEqual.Plugin,
	addtend.Plugin,
	extractresult.Plugin,
	removeunusedimport.Plugin,
	removeunusedvariable.Plugin,
}
