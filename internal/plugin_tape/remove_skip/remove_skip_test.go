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

	Test(t, "remove-skip: report: skip-skip", func(t *T) {
		t.Report("skip-skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform: skip-skip", func(t *T) {
		t.Transform("skip-skip")
		t.End()
	})

	Test(t, "remove-skip: report: with-options", func(t *T) {
		t.Report("with-options", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform: with-options", func(t *T) {
		t.Transform("with-options")
		t.End()
	})

	Test(t, "remove-skip: no report: not-test", func(t *T) {
		t.NoReport("not-test")
		t.End()
	})

	Test(t, "remove-skip: no transform: not-test", func(t *T) {
		t.NoTransform("not-test")
		t.End()
	})
}
