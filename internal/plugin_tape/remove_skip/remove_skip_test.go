package remove_skip_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_skip"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-skip", remove_skip.Plugin{})

func TestRemoveSkip(t *testing.T) {
	Test(t, "remove-skip: report: remove-skip", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform: remove-skip", func(t *T) {
		t.Transform("remove-skip")
		t.End()
	})

	Test(t, "remove-skip: no report: no-skip", func(t *T) {
		t.NoReport("no-skip")
		t.End()
	})

	Test(t, "remove-skip: no transform: no-skip", func(t *T) {
		t.NoTransform("no-skip")
		t.End()
	})
}
