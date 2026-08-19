package convert_equal_to_deep_equal_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_deep_equal"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-deep-equal", convert_equal_to_deep_equal.Plugin{})

func TestConvertEqualToDeepEqual(t *testing.T) {
	Test(t, "convert-equal-to-deep-equal: report: equal-slice", func(t *T) {
		t.Report("equal-slice", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform: equal-slice", func(t *T) {
		t.Transform("equal-slice")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report: equal-slice-first", func(t *T) {
		t.Report("equal-slice-first", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report: convert-equal-to-deep-equal", func(t *T) {
		t.Report("convert-equal-to-deep-equal", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform: convert-equal-to-deep-equal", func(t *T) {
		t.Transform("convert-equal-to-deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no report: deep-equal", func(t *T) {
		t.NoReport("deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no transform: deep-equal", func(t *T) {
		t.NoTransform("deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report: not-equal-slice", func(t *T) {
		t.Report("not-equal-slice", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform: not-equal-slice", func(t *T) {
		t.Transform("not-equal-slice")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report: named-var", func(t *T) {
		t.Report("named-var", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform: named-var", func(t *T) {
		t.Transform("named-var")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no report: named-var-no-report", func(t *T) {
		t.NoReport("named-var-no-report")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no transform: named-var-no-report", func(t *T) {
		t.NoTransform("named-var-no-report")
		t.End()
	})
}
