package removeskip_test

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
	Test           = indratest.CreateTest(Plugin, _dir)
)

func TestRemoveSkip(t *testing.T) {
	Test(t, "remove-skip: report skip call", func(t *indratest.T) {
		t.Report("skip", "remove t.Skip call")
		t.End()
	})

	Test(t, "remove-skip: no report for clean file", func(t *indratest.T) {
		t.NoReport("no-skip")
		t.End()
	})
}
