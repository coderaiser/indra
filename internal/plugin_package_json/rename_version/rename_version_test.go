package rename_version_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_package_json/rename_version"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("package-json/rename-version", rename_version.Plugin{})

func TestRenameVersion(t *testing.T) {
	Test(t, "package-json/rename-version: report: rename-version", func(t *T) {
		t.Report("rename-version", "normalise version field")
		t.End()
	})

	Test(t, "package-json/rename-version: transform: rename-version", func(t *T) {
		t.Transform("rename-version")
		t.End()
	})

	Test(t, "package-json/rename-version: no report: no-rename", func(t *T) {
		t.NoReport("no-rename")
		t.End()
	})

	Test(t, "package-json/rename-version: no transform: no-rename", func(t *T) {
		t.NoTransform("no-rename")
		t.End()
	})
}
