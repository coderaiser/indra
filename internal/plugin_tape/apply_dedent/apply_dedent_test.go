package apply_dedent_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/apply_dedent"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("tape/apply-dedent", apply_dedent.Plugin{})

func TestApplyDedent(t *testing.T) {
	Test(t, "tape/apply-dedent: report: with-dedent", func(t *T) {
		t.Report("with-dedent", "apply dedent")
		t.End()
	})

	Test(t, "tape/apply-dedent: transform: with-dedent", func(t *T) {
		t.Transform("with-dedent")
		t.End()
	})

	Test(t, "tape/apply-dedent: no report: no-dedent", func(t *T) {
		t.NoReport("no-dedent")
		t.End()
	})

	Test(t, "tape/apply-dedent: no report: no-match", func(t *T) {
		t.NoReport("no-match")
		t.End()
	})
}
