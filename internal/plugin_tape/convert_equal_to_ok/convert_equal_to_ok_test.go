package convert_equal_to_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-ok", convert_equal_to_ok.Plugin{})

func TestConvertEqualToOk(t *testing.T) {
	Test(t, "convert-equal-to-ok: report Equal(x, true) convert-equal-to-ok", func(t *T) {
		t.Report("convert-equal-to-ok", "convert Equal(x, true) to Ok(x)")
		t.End()
	})

	Test(t, "convert-equal-to-ok: transform Equal(x, true) to Ok convert-equal-to-ok", func(t *T) {
		t.Transform("convert-equal-to-ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no report for Ok ok", func(t *T) {
		t.NoReport("ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no transform for Ok ok", func(t *T) {
		t.NoTransform("ok")
		t.End()
	})
}
