package remove_unused_import_test

import (
	"path/filepath"
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var (
	_, _file, _, _ = runtime.Caller(0)
	_dir           = filepath.Join(filepath.Dir(_file), "fixture")
	Test           = indratest.CreateTest("coderaiser/indra/internal/plugins/remove-unused-import", _dir)
)

func TestRemoveUnusedImport(t *testing.T) {
	Test(t, "remove-unused-import: report unused import", func(t *indratest.T) {
		t.Report("unused-import", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-import: no report when all imports used", func(t *indratest.T) {
		t.NoReport("used-import")
		t.End()
	})

	Test(t, "remove-unused-import: no report for blank import", func(t *indratest.T) {
		t.NoReport("blank-import")
		t.End()
	})

	Test(t, "remove-unused-import: no report for dot import", func(t *indratest.T) {
		t.NoReport("dot-import")
		t.End()
	})

	Test(t, "remove-unused-import: report unused aliased import", func(t *indratest.T) {
		t.Report("alias-unused", `remove unused import: "fmt"`)
		t.End()
	})

	Test(t, "remove-unused-import: no report for used alias", func(t *indratest.T) {
		t.NoReport("alias-used")
		t.End()
	})

	Test(t, "remove-unused-import: fix removes unused import", func(t *indratest.T) {
		t.Transform("unused-import")
		t.End()
	})

	Test(t, "remove-unused-import: fix removes unused aliased import", func(t *indratest.T) {
		t.Transform("alias-unused")
		t.End()
	})

	Test(t, "remove-unused-import: fix keeps used import in mixed block", func(t *indratest.T) {
		t.Transform("mixed-imports")
		t.End()
	})
}
