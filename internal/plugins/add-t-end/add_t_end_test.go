package addtend_test

import (
	"path/filepath"
	"runtime"
	"testing"

	. "coderaiser/indra/internal/plugins/add-t-end"
	indratest "coderaiser/indra/internal/test"
)

var (
	_, _file, _, _ = runtime.Caller(0)
	_dir           = filepath.Join(filepath.Dir(_file), "fixture")
	Test           = indratest.CreateTest(Plugin, _dir)
)

func TestAddTEnd(t *testing.T) {
	Test(t, "add-t-end: report missing End in tape.Test", func(t *indratest.T) {
		t.Report("missing-end", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform adds End to tape.Test", func(t *indratest.T) {
		t.Transform("missing-end")
		t.End()
	})

	Test(t, "add-t-end: no report when End already present", func(t *indratest.T) {
		t.NoReport("has-end")
		t.End()
	})

	Test(t, "add-t-end: no transform when End already present", func(t *indratest.T) {
		t.NoTransform("has-end")
		t.End()
	})

	Test(t, "add-t-end: report missing End in tape.Only", func(t *indratest.T) {
		t.Report("missing-end-only", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform adds End to tape.Only", func(t *indratest.T) {
		t.Transform("missing-end-only")
		t.End()
	})
}
