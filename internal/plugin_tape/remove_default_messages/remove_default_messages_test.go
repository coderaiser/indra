package remove_default_messages_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_default_messages"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-default-messages", remove_default_messages.Plugin{})

func TestRemoveDefaultMessages(t *testing.T) {
	Test(t, "remove-default-messages: report: remove-default-messages", func(t *T) {
		t.Report("remove-default-messages", "Avoid passing default messages to operators")
		t.End()
	})

	Test(t, "remove-default-messages: transform: remove-default-messages", func(t *T) {
		t.Transform("remove-default-messages")
		t.End()
	})

	Test(t, "remove-default-messages: no report: custom-message", func(t *T) {
		t.NoReport("custom-message")
		t.End()
	})

	Test(t, "remove-default-messages: no transform: custom-message", func(t *T) {
		t.NoTransform("custom-message")
		t.End()
	})
}
