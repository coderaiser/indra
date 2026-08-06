package convert_equal_to_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_ok"
	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.For("convert-equal-to-ok", convert_equal_to_ok.Plugin{})

func TestConvertEqualToOk(t *testing.T) {
	Test(t, "convert-equal-to-ok: report Equal(x, true)", func(t *indratest.T) {
		t.Report("convert-equal-to-ok", "convert Equal(x, true) to Ok(x)")
		t.End()
	})

	Test(t, "convert-equal-to-ok: transform Equal(x, true) to Ok", func(t *indratest.T) {
		t.Transform("convert-equal-to-ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no report for Ok", func(t *indratest.T) {
		t.NoReport("ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no transform for Ok", func(t *indratest.T) {
		t.NoTransform("ok")
		t.End()
	})
}
