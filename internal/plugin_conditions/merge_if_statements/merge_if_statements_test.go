package merge_if_statements_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/merge_if_statements"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/merge-if-statements", merge_if_statements.Plugin{})

func TestMergeIfStatements(t *testing.T) {
	Test(t, "conditions/merge-if-statements: report: report: merge-if-statements", func(t *T) {
		t.Report("merge-if-statements", "merge if statements")
		t.End()
	})

	Test(t, "conditions/merge-if-statements: transform: transform: merge-if-statements", func(t *T) {
		t.Transform("merge-if-statements")
		t.End()
	})

	Test(t, "conditions/merge-if-statements: no report: no report: no report: no-merge-outer-else", func(t *T) {
		t.NoReport("no-merge-outer-else")
		t.End()
	})

	Test(t, "conditions/merge-if-statements: no report: no report: no report: no-merge-inner-else", func(t *T) {
		t.NoReport("no-merge-inner-else")
		t.End()
	})

	Test(t, "conditions/merge-if-statements: no report: no report: no report: no-merge-multi-stmt", func(t *T) {
		t.NoReport("no-merge-multi-stmt")
		t.End()
	})

	Test(t, "conditions/merge-if-statements: no report: no report: no report: no-merge-inner-init", func(t *T) {
		t.NoReport("no-merge-inner-init")
		t.End()
	})
}
