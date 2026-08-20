package convert_equal_to_not_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_not_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-not-ok", convert_equal_to_not_ok.Plugin{})

func TestConvertEqualToNotOk(t *testing.T) {
	Test(t, "convert-equal-to-not-ok: report: convert-equal-to-not-ok", func(t *T) {
		t.Report("convert-equal-to-not-ok", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: convert-equal-to-not-ok", func(t *T) {
		t.Transform("convert-equal-to-not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: not-ok", func(t *T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: not-ok", func(t *T) {
		t.NoTransform("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report: convert-equal-false-to-not-ok", func(t *T) {
		t.Report("convert-equal-false-to-not-ok", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: convert-equal-false-to-not-ok", func(t *T) {
		t.Transform("convert-equal-false-to-not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: empty-string", func(t *T) {
		t.Transform("empty-string")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report: zero-int", func(t *T) {
		t.Report("zero-int", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: zero-int", func(t *T) {
		t.Transform("zero-int")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: zero-float", func(t *T) {
		t.NoReport("zero-float")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: zero-float", func(t *T) {
		t.NoTransform("zero-float")
		t.End()
	})
}
