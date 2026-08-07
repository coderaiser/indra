package convert_no_error_to_not_ok_test

import (
	"go/ast"
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_no_error_to_not_ok"
	"coderaiser/indra/types"
	. "coderaiser/indra/internal/test"

	tape "github.com/coderaiser/go-tape"
)

var Test = CreateTest("convert-no-error-to-not-ok", convert_no_error_to_not_ok.Plugin{})

// TestFixDirect covers the Fix early return for a file without the go-tape
// import, which is not reachable through the push-based fixture harness.
func TestFixDirect(t *testing.T) {
	tape.Test(t, "fix: no-op without go-tape import", func(t *tape.T) {
		convert_no_error_to_not_ok.Fix(types.Path{Node: &ast.File{}}, nil)
		t.Pass("returned without panic")
		t.End()
	})
}

func TestConvertNoErrorToNotOk(t *testing.T) {
	Test(t, "convert-no-error-to-not-ok: reports NoError call", func(t *T) {
		t.Report("convert-no-error-to-not-ok", "convert NoError(err) to NotOk(err)")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: fixes NoError to NotOk", func(t *T) {
		t.Transform("convert-no-error-to-not-ok")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no report for NotOk", func(t *T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no report without go-tape import", func(t *T) {
		t.NoReport("no-error-other-lib")
		t.End()
	})

	Test(t, "convert-no-error-to-not-ok: no transform without go-tape import", func(t *T) {
		t.NoTransform("no-error-other-lib")
		t.End()
	})
}
