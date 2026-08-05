package convert_no_error_to_not_ok_test

import (
	"go/ast"
	"runtime"
	"testing"

	"coderaiser/indra/internal/plugins/convert_no_error_to_not_ok"
	indratest "coderaiser/indra/internal/test"

	tape "github.com/coderaiser/go-tape"
)

var Test = indratest.CreateTest(runtime.Caller(0))

// TestFixDirect covers the Fix early return for a file without the go-tape
// import, which is not reachable through the push-based fixture harness.
func TestFixDirect(t *testing.T) {
	tape.Test(t, "fix: no-op without go-tape import", func(t *tape.T) {
		convert_no_error_to_not_ok.Fix(&ast.File{}, nil)
		t.Pass("returned without panic")
		t.End()
	})
}

func TestConvertNoErrorToNotOk(t *testing.T) {
	Test(t, "convert-no-error-to-not-ok: reports NoError call", func(t *indratest.T) {
		t.Report("convert-no-error-to-not-ok", "convert NoError(err) to NotOk(err)")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: fixes NoError to NotOk", func(t *indratest.T) {
		t.Transform("convert-no-error-to-not-ok")
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
