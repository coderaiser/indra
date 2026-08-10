package remove_useless_match_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/remove_useless_match"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-match", remove_useless_match.Plugin{})

func TestRemoveUselessMatch(t *testing.T) {
	Test(t, "remove-useless-match: report: useless-nil", func(t *T) {
		t.Report("useless-nil", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: report: useless-empty", func(t *T) {
		t.Report("useless-empty", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: report: remove-useless-match", func(t *T) {
		t.Report("remove-useless-match", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: no report: useful-match", func(t *T) {
		t.NoReport("useful-match")
		t.End()
	})

	Test(t, "remove-useless-match: transform: useless-nil", func(t *T) {
		t.Transform("useless-nil")
		t.End()
	})

	Test(t, "remove-useless-match: transform: useless-empty", func(t *T) {
		t.Transform("useless-empty")
		t.End()
	})

	Test(t, "remove-useless-match: no report: mixed-guards", func(t *T) {
		t.NoReport("mixed-guards")
		t.End()
	})

	Test(t, "remove-useless-match: no transform: mixed-guards", func(t *T) {
		t.NoTransform("mixed-guards")
		t.End()
	})

	Test(t, "remove-useless-match: no report: wrong-return-type", func(t *T) {
		t.NoReport("wrong-return-type")
		t.End()
	})

	Test(t, "remove-useless-match: no transform: wrong-return-type", func(t *T) {
		t.NoTransform("wrong-return-type")
		t.End()
	})

	Test(t, "remove-useless-match: no report: non-composite-return", func(t *T) {
		t.NoReport("non-composite-return")
		t.End()
	})

	Test(t, "remove-useless-match: no transform: non-composite-return", func(t *T) {
		t.NoTransform("non-composite-return")
		t.End()
	})

	Test(t, "remove-useless-match: no report: no-results", func(t *T) {
		t.NoReport("no-results")
		t.End()
	})

	Test(t, "remove-useless-match: no report: non-return-body", func(t *T) {
		t.NoReport("non-return-body")
		t.End()
	})

	Test(t, "remove-useless-match: no report: ident-return", func(t *T) {
		t.NoReport("ident-return")
		t.End()
	})
}
