package convert_equal_to_not_ok_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestConvertEqualToNotOk(t *testing.T) {
	Test(t, "convert-equal-to-not-ok: report Equal(err, nil)", func(t *indratest.T) {
		t.Report("convert-equal-to-not-ok", "convert Equal(err, nil) to NotOk(err)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform Equal(err, nil) to NotOk", func(t *indratest.T) {
		t.Transform("convert-equal-to-not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report for NotOk", func(t *indratest.T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform for NotOk", func(t *indratest.T) {
		t.NoTransform("not-ok")
		t.End()
	})
}
