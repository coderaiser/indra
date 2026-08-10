package convert_no_error_to_not_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_no_error_to_not_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-no-error-to-not-ok", convert_no_error_to_not_ok.Plugin{})

func TestConvertNoErrorToNotOk(t *testing.T) {
	Test(t, "convert-no-error-to-not-ok: report: convert-no-error-to-not-ok", func(t *T) {
		t.Report("convert-no-error-to-not-ok", "convert NoError(err) to NotOk(err)")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: transform: convert-no-error-to-not-ok", func(t *T) {
		t.Transform("convert-no-error-to-not-ok")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no report: not-ok", func(t *T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no report: no-error-other-lib", func(t *T) {
		t.NoReport("no-error-other-lib")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no transform: no-error-other-lib", func(t *T) {
		t.NoTransform("no-error-other-lib")
		t.End()
	})
}
