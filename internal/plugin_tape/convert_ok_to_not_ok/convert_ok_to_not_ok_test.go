package convert_ok_to_not_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_ok_to_not_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-ok-to-not-ok", convert_ok_to_not_ok.Plugin{})

func TestConvertOkToNotOk(t *testing.T) {
	Test(t, "convert-ok-to-not-ok: report Ok(err == nil) convert-ok-to-not-ok", func(t *T) {
		t.Report("convert-ok-to-not-ok", "convert Ok to NotOk")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: transform Ok(err == nil) to NotOk convert-ok-to-not-ok", func(t *T) {
		t.Transform("convert-ok-to-not-ok")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: report Ok(!err) ok-not", func(t *T) {
		t.Report("ok-not", "convert Ok to NotOk")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: transform Ok(!err) to NotOk ok-not", func(t *T) {
		t.Transform("ok-not")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: no report for NotOk not-ok", func(t *T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: no transform for NotOk not-ok", func(t *T) {
		t.NoTransform("not-ok")
		t.End()
	})
}
