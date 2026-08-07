package remove_useless_match_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/remove_useless_match"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-match", remove_useless_match.Plugin{})

func TestRemoveUselessMatch(t *testing.T) {
	Test(t, "remove-useless-match: report nil guard entry useless-nil", func(t *T) {
		t.Report("useless-nil", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: report empty Matcher useless-empty", func(t *T) {
		t.Report("useless-empty", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: report useless match remove-useless-match", func(t *T) {
		t.Report("remove-useless-match", "remove useless Match")
		t.End()
	})

	Test(t, "remove-useless-match: no report for real guard useful-match", func(t *T) {
		t.NoReport("useful-match")
		t.End()
	})

	Test(t, "remove-useless-match: transform nil guard entry useless-nil", func(t *T) {
		t.Transform("useless-nil")
		t.End()
	})

	Test(t, "remove-useless-match: transform empty Matcher useless-empty", func(t *T) {
		t.Transform("useless-empty")
		t.End()
	})
}
