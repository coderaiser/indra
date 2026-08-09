package remove_useless_comments_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/remove_useless_comments"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/remove-useless-comments", remove_useless_comments.Plugin{})

func TestRemoveUselessComments(t *testing.T) {
	Test(t, "conditions/remove-useless-comments: report: remove-useless-comments", func(t *T) {
		t.Report("remove-useless-comments", "remove useless comments")
		t.End()
	})

	Test(t, "conditions/remove-useless-comments: transform: remove-useless-comments", func(t *T) {
		t.Transform("remove-useless-comments")
		t.End()
	})

	Test(t, "conditions/remove-useless-comments: no report: no-separator", func(t *T) {
		t.NoReport("no-separator")
		t.End()
	})

	Test(t, "conditions/remove-useless-comments: no report: short-dash", func(t *T) {
		t.NoReport("short-dash")
		t.End()
	})
}
