package remove_skip_test

import (
	"path/filepath"
	"runtime"
	"testing"

	. "coderaiser/indra/internal/plugins/remove-skip"
	indratest "coderaiser/indra/internal/test"
)

var (
	_, _file, _, _ = runtime.Caller(0)
	_dir           = filepath.Join(filepath.Dir(_file), "fixture")
	Test           = indratest.CreateTest(Self, _dir)
)

func TestRemoveSkip(t *testing.T) {
	Test(t, "remove-skip: report Test.Skip call", func(t *indratest.T) {
		t.Report("skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform Test.Skip to Test", func(t *indratest.T) {
		t.Transform("skip")
		t.End()
	})

	Test(t, "remove-skip: no report for Test call", func(t *indratest.T) {
		t.NoReport("no-skip")
		t.End()
	})

	Test(t, "remove-skip: no transform for Test call", func(t *indratest.T) {
		t.NoTransform("no-skip")
		t.End()
	})
}
