package convert_ok_to_not_ok_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestConvertOkToNotOk(t *testing.T) {
	Test(t, "convert-ok-to-not-ok: report Ok(err == nil)", func(t *indratest.T) {
		t.Report("convert-ok-to-not-ok", "convert Ok(err == nil) to NotOk(err)")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: transform Ok(err == nil) to NotOk", func(t *indratest.T) {
		t.Transform("convert-ok-to-not-ok")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: no report for NotOk", func(t *indratest.T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-ok-to-not-ok: no transform for NotOk", func(t *indratest.T) {
		t.NoTransform("not-ok")
		t.End()
	})
}
