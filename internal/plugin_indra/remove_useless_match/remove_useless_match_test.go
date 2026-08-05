package remove_useless_match_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/remove_useless_match"
	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.For("remove-useless-match", remove_useless_match.Plugin{})

func TestRemoveUselessMatch(t *testing.T) {
	Test(t, "remove-useless-match: report nil guard entry", func(t *indratest.T) {
		t.Report("useless-nil", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: report empty Matcher", func(t *indratest.T) {
		t.Report("useless-empty", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: report useless match", func(t *indratest.T) {
		t.Report("remove-useless-match", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: no report for real guard", func(t *indratest.T) {
		t.NoReport("useful-match")
		t.End()
	})

	Test(t, "remove-useless-match: transform nil guard entry", func(t *indratest.T) {
		t.Transform("useless-nil")
		t.End()
	})

	Test(t, "remove-useless-match: transform empty Matcher", func(t *indratest.T) {
		t.Transform("useless-empty")
		t.End()
	})
}
