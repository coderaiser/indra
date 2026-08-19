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

	Test(t, "add-t-end: transform: missing-end", func(t *T) {
		t.Transform("missing-end")
		t.End()
	})

	Test(t, "add-t-end: report: missing-end-only", func(t *T) {
		t.Report("missing-end-only", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform: missing-end-only", func(t *T) {
		t.Transform("missing-end-only")
		t.End()
	})

	Test(t, "add-t-end: report: add-t-end", func(t *T) {
		t.Report("add-t-end", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform: add-t-end", func(t *T) {
		t.Transform("add-t-end")
		t.End()
	})

	Test(t, "add-t-end: no report: has-end", func(t *T) {
		t.NoReport("has-end")
		t.End()
	})

	Test(t, "add-t-end: no transform: has-end", func(t *T) {
		t.NoTransform("has-end")
		t.End()
	})

	Test(t, "add-t-end: no report: assign-end", func(t *T) {
		t.NoReport("assign-end")
		t.End()
	})

	Test(t, "add-t-end: no transform: assign-end", func(t *T) {
		t.NoTransform("assign-end")
		t.End()
	})

	Test(t, "add-t-end: no report: contains-end", func(t *T) {
		t.NoReport("contains-end")
		t.End()
	})

	Test(t, "add-t-end: no transform: contains-end", func(t *T) {
		t.NoTransform("contains-end")
		t.End()
	})

	Test(t, "add-t-end: no report: callback-end", func(t *T) {
		t.NoReport("callback-end")
		t.End()
	})

	Test(t, "add-t-end: no transform: callback-end", func(t *T) {
		t.NoTransform("callback-end")
		t.End()
	})

	Test(t, "add-t-end: report: empty-body", func(t *T) {
		t.Report("empty-body", "tape: missing t.End()")
		t.End()
	})

	Test(t, "add-t-end: transform: empty-body", func(t *T) {
		t.Transform("empty-body")
		t.End()
	})
}
