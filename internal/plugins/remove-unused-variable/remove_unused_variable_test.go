package remove_unused_variable_test

import (
	"path/filepath"
	"runtime"
	"testing"

	. "coderaiser/indra/internal/plugins/remove-unused-variable"
	indratest "coderaiser/indra/internal/test"
)

var (
	_, _file, _, _ = runtime.Caller(0)
	_dir           = filepath.Join(filepath.Dir(_file), "fixture")
	Test           = indratest.CreateTest(Plugin, _dir)
)

func TestRemoveUnusedVariable(t *testing.T) {
	Test(t, "remove-unused-variable: report unused var", func(t *indratest.T) {
		t.Report("unused-var", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variable: no report when all vars used", func(t *indratest.T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variable: no report for blank in tuple assign", func(t *indratest.T) {
		t.NoReport("tuple-blank")
		t.End()
	})
}
