package replace_test_message_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/replace_test_message"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("replace-test-message", replace_test_message.Plugin{})

func TestReplaceTestMessage(t *testing.T) {
	Test(t, "replace-test-message: report replace-test-message", func(t *T) {
		t.Report("replace-test-message", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform replace-test-message", func(t *T) {
		t.Transform("replace-test-message")
		t.End()
	})

	Test(t, "replace-test-message: no report already-correct", func(t *T) {
		t.NoReport("already-correct")
		t.End()
	})

	Test(t, "replace-test-message: no report no-fixture-call", func(t *T) {
		t.NoReport("no-fixture-call")
		t.End()
	})
}
