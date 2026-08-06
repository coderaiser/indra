package convert_equal_to_not_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_not_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-not-ok", convert_equal_to_not_ok.Plugin{})

func TestConvertEqualToNotOk(t *testing.T) {
	Test(t, "convert-equal-to-not-ok: report Equal(err, nil)", func(t *T) {
		t.Report("convert-equal-to-not-ok", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform Equal(err, nil) to NotOk", func(t *T) {
		t.Transform("convert-equal-to-not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report for NotOk", func(t *T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform for NotOk", func(t *T) {
		t.NoTransform("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report Equal(x, false)", func(t *T) {
		t.Report("convert-equal-false-to-not-ok", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform Equal(x, false) to NotOk", func(t *T) {
		t.Transform("convert-equal-false-to-not-ok")
		t.End()
	})
}
