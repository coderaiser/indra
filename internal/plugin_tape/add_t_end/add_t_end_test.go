package add_t_end_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/add_t_end"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("add-t-end", add_t_end.Plugin{})

func TestAddTEnd(t *testing.T) {
	Test(t, "add-t-end: report: missing-end", func(t *T) {
		t.Report("missing-end", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform adds End to Test", func(t *T) {
		t.Transform("missing-end")
		t.End()
	})

	Test(t, "add-t-end: report missing End in Test.Only", func(t *T) {
		t.Report("missing-end-only", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform adds End to Test.Only", func(t *T) {
		t.Transform("missing-end-only")
		t.End()
	})

	Test(t, "add-t-end: report missing End in Test", func(t *T) {
		t.Report("add-t-end", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform: add-t-end", func(t *T) {
		t.Transform("add-t-end")
		t.End()
	})

	Test(t, "add-t-end: no report when End present", func(t *T) {
		t.NoReport("has-end")
		t.End()
	})

	Test(t, "add-t-end: no transform when End present", func(t *T) {
		t.NoTransform("has-end")
		t.End()
	})
}
