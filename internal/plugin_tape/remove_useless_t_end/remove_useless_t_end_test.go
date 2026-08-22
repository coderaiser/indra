package remove_useless_t_end_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_useless_t_end"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-t-end", remove_useless_t_end.Plugin{})

func TestRemoveUselessTEnd(t *testing.T) {
	Test(t, "remove-useless-t-end: report: remove-useless-t-end", func(t *T) {
		t.Report("remove-useless-t-end", `Avoid useless "t.End()"`)
		t.End()
	})

	Test(t, "remove-useless-t-end: transform: remove-useless-t-end", func(t *T) {
		t.Transform("remove-useless-t-end")
		t.End()
	})

	Test(t, "remove-useless-t-end: no report: one-end", func(t *T) {
		t.NoReport("one-end")
		t.End()
	})

	Test(t, "remove-useless-t-end: no transform: one-end", func(t *T) {
		t.NoTransform("one-end")
		t.End()
	})

	Test(t, "remove-useless-t-end: no report: no-end", func(t *T) {
		t.NoReport("no-end")
		t.End()
	})

	Test(t, "remove-useless-t-end: no transform: no-end", func(t *T) {
		t.NoTransform("no-end")
		t.End()
	})
}
