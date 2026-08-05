package convert_no_error_to_not_ok_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestConvertNoErrorToNotOk(t *testing.T) {
	Test(t, "convert-no-error-to-not-ok: reports NoError call", func(t *indratest.T) {
		t.Report("no-error", "convert NoError(err) to NotOk(err)")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: fixes NoError to NotOk", func(t *indratest.T) {
		t.Transform("no-error")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no report for NotOk", func(t *indratest.T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no report without go-tape import", func(t *indratest.T) {
		t.NoReport("no-error-other-lib")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no transform without go-tape import", func(t *indratest.T) {
		t.NoTransform("no-error-other-lib")
		t.End()
	})
}
