package convert_equal_to_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-ok", convert_equal_to_ok.Plugin{})

func TestConvertEqualToOk(t *testing.T) {
	Test(t, "convert-equal-to-ok: report: convert-equal-to-ok", func(t *T) {
		t.Report("convert-equal-to-ok", "convert Equal(x, true) to Ok(x)")
		t.End()
	})

	Test(t, "convert-equal-to-ok: transform: convert-equal-to-ok", func(t *T) {
		t.Transform("convert-equal-to-ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no report: ok", func(t *T) {
		t.NoReport("ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no transform: ok", func(t *T) {
		t.NoTransform("ok")
		t.End()
	})

	Test(t, "convert-equal-to-ok: report: message", func(t *T) {
		t.Report("message", "convert Equal(x, true) to Ok(x)")
		t.End()
	})

	Test(t, "convert-equal-to-ok: transform: message", func(t *T) {
		t.Transform("message")
		t.End()
	})

	Test(t, "convert-equal-to-ok: report: deep-equal", func(t *T) {
		t.Report("deep-equal", "convert Equal(x, true) to Ok(x)")
		t.End()
	})

	Test(t, "convert-equal-to-ok: transform: deep-equal", func(t *T) {
		t.Transform("deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no report: false", func(t *T) {
		t.NoReport("false")
		t.End()
	})

	Test(t, "convert-equal-to-ok: no transform: false", func(t *T) {
		t.NoTransform("false")
		t.End()
	})

	Test(t, "convert-equal-to-ok: report: deep-equal-message", func(t *T) {
		t.Report("deep-equal-message", "convert Equal(x, true) to Ok(x)")
		t.End()
	})

	Test(t, "convert-equal-to-ok: transform: deep-equal-message", func(t *T) {
		t.Transform("deep-equal-message")
		t.End()
	})
}
