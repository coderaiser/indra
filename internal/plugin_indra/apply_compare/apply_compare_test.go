package apply_compare_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/apply_compare"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("apply-compare", apply_compare.Plugin{})

func TestApplyCompare(t *testing.T) {
	Test(t, "apply-compare: report: apply-compare", func(t *T) {
		t.Report("apply-compare", "use Compare instead of GetTemplateValues != nil")
		t.End()
	})

	Test(t, "apply-compare: transform: apply-compare", func(t *T) {
		t.Transform("apply-compare")
		t.End()
	})

	Test(t, "apply-compare: no report: no-plugin", func(t *T) {
		t.NoReport("no-plugin")
		t.End()
	})

	Test(t, "apply-compare: no report: no-match", func(t *T) {
		t.NoReport("no-match")
		t.End()
	})
}
