package remove_skip_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_skip"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-skip", remove_skip.Plugin{})

func TestRemoveSkip(t *testing.T) {
	Test(t, "remove-skip: report Test.Skip call remove-skip", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform Test.Skip to Test remove-skip", func(t *T) {
		t.Transform("remove-skip")
		t.End()
	})

	Test(t, "remove-skip: no report for Test call no-skip", func(t *T) {
		t.NoReport("no-skip")
		t.End()
	})

	Test(t, "remove-skip: no transform for Test call no-skip", func(t *T) {
		t.NoTransform("no-skip")
		t.End()
	})
}
