package merge_if_with_else_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/merge_if_with_else"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/merge-if-with-else", merge_if_with_else.Plugin{})

func TestMergeIfWithElse(t *testing.T) {
	Test(t, "conditions/merge-if-with-else: report: report: merge-if-with-else", func(t *T) {
		t.Report("merge-if-with-else", "merge if with else")
		t.End()
	})

	Test(t, "conditions/merge-if-with-else: transform: transform: merge-if-with-else", func(t *T) {
		t.Transform("merge-if-with-else")
		t.End()
	})

	Test(t, "conditions/merge-if-with-else: no report: no report: no report: no-merge-different-bodies", func(t *T) {
		t.NoReport("no-merge-different-bodies")
		t.End()
	})

	Test(t, "conditions/merge-if-with-else: no report: no report: no report: no-merge-no-else", func(t *T) {
		t.NoReport("no-merge-no-else")
		t.End()
	})

	Test(t, "conditions/merge-if-with-else: no report: no report: no report: no-merge-plain-else", func(t *T) {
		t.NoReport("no-merge-plain-else")
		t.End()
	})

	Test(t, "conditions/merge-if-with-else: no report: no report: no report: no-merge-else-chain", func(t *T) {
		t.NoReport("no-merge-else-chain")
		t.End()
	})

	Test(t, "conditions/merge-if-with-else: no report: no report: no report: no-merge-init", func(t *T) {
		t.NoReport("no-merge-init")
		t.End()
	})
}
