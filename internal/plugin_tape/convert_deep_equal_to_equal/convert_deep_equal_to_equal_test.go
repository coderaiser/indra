package convert_deep_equal_to_equal_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_deep_equal_to_equal"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-deep-equal-to-equal", convert_deep_equal_to_equal.Plugin{})

func TestConvertDeepEqualToEqual(t *testing.T) {
	Test(t, "convert-deep-equal-to-equal: report: convert-deep-equal-to-equal", func(t *T) {
		t.Report("convert-deep-equal-to-equal", "Use Equal() when comparing primitives")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: transform: convert-deep-equal-to-equal", func(t *T) {
		t.Transform("convert-deep-equal-to-equal")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: transform: nil-arg", func(t *T) {
		t.Transform("nil-arg")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: transform: bool-arg", func(t *T) {
		t.Transform("bool-arg")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: transform: int-arg", func(t *T) {
		t.Transform("int-arg")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: report: message", func(t *T) {
		t.Report("message", "Use Equal() when comparing primitives")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: transform: message", func(t *T) {
		t.Transform("message")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: no report: slice-arg", func(t *T) {
		t.NoReport("slice-arg")
		t.End()
	})

	Test(t, "convert-deep-equal-to-equal: no transform: slice-arg", func(t *T) {
		t.NoTransform("slice-arg")
		t.End()
	})
}
