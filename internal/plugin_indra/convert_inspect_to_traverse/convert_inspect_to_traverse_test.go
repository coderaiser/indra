package convert_inspect_to_traverse_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/convert_inspect_to_traverse"
	. "coderaiser/indra/internal/test"
	"coderaiser/indra/types"
)

var Test = CreateTest("convert-inspect-to-traverse", convert_inspect_to_traverse.Plugin{})

func TestConvertInspectToTraverse(t *testing.T) {
	Test(t, "convert-inspect-to-traverse: report inspect", func(t *T) {
		t.Report("inspect", "convert ast.Inspect to path.Traverse")
		t.End()
	})

	Test(t, "convert-inspect-to-traverse: no report no-plugin", func(t *T) {
		t.NoReport("no-plugin")
		t.End()
	})

	Test(t, "convert-inspect-to-traverse: Fix is a no-op", func(t *T) {
		convert_inspect_to_traverse.Fix(types.Path{}, nil)
		convert_inspect_to_traverse.Plugin{}.Fix(types.Path{}, nil)
		t.Pass("no panic")
		t.End()
	})
}
