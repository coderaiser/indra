package remove_unused_import_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

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
