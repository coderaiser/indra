package convert_equal_to_deep_equal_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_deep_equal"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-deep-equal", convert_equal_to_deep_equal.Plugin{})

func TestConvertEqualToDeepEqual(t *testing.T) {
	Test(t, "convert-equal-to-deep-equal: report array second arg", func(t *T) {
		t.Report("equal-slice", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform array second arg", func(t *T) {
		t.Transform("equal-slice")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report array first arg", func(t *T) {
		t.Report("equal-slice-first", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report Equal with slice", func(t *T) {
		t.Report("convert-equal-to-deep-equal", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform Equal with slice", func(t *T) {
		t.Transform("convert-equal-to-deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no report for DeepEqual", func(t *T) {
		t.NoReport("deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no transform for DeepEqual", func(t *T) {
		t.NoTransform("deep-equal")
		t.End()
	})
}
