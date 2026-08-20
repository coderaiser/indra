package replace_test_message_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/replace_test_message"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("replace-test-message", replace_test_message.Plugin{})

func TestReplaceTestMessage(t *testing.T) {
	Test(t, "replace-test-message: report: replace-test-message", func(t *T) {
		t.Report("replace-test-message", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: replace-test-message", func(t *T) {
		t.Transform("replace-test-message")
		t.End()
	})

	Test(t, "replace-test-message: no report: already-correct", func(t *T) {
		t.NoReport("already-correct")
		t.End()
	})

	Test(t, "replace-test-message: no report: no-fixture-call", func(t *T) {
		t.NoReport("no-fixture-call")
		t.End()
	})

	Test(t, "replace-test-message: report: no-separator", func(t *T) {
		t.Report("no-separator", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: no-separator", func(t *T) {
		t.Transform("no-separator")
		t.End()
	})

	Test(t, "replace-test-message: no report: odd-callback", func(t *T) {
		t.NoReport("odd-callback")
		t.End()
	})

	Test(t, "replace-test-message: no transform: odd-callback", func(t *T) {
		t.NoTransform("odd-callback")
		t.End()
	})

	Test(t, "replace-test-message: no report: non-literal-call", func(t *T) {
		t.NoReport("non-literal-call")
		t.End()
	})

	Test(t, "replace-test-message: no transform: non-literal-call", func(t *T) {
		t.NoTransform("non-literal-call")
		t.End()
	})

	Test(t, "replace-test-message: report: mixed-args", func(t *T) {
		t.Report("mixed-args", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: mixed-args", func(t *T) {
		t.Transform("mixed-args")
		t.End()
	})

	Test(t, "replace-test-message: report: no-report-verb", func(t *T) {
		t.Report("no-report-verb", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: no-report-verb", func(t *T) {
		t.Transform("no-report-verb")
		t.End()
	})

	Test(t, "replace-test-message: report: no-transform-verb", func(t *T) {
		t.Report("no-transform-verb", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: no-transform-verb", func(t *T) {
		t.Transform("no-transform-verb")
		t.End()
	})

	Test(t, "replace-test-message: report: wrong-verb", func(t *T) {
		t.Report("wrong-verb", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: wrong-verb", func(t *T) {
		t.Transform("wrong-verb")
		t.End()
	})

	Test(t, "replace-test-message: report: wrong-fixture", func(t *T) {
		t.Report("wrong-fixture", "replace test message")
		t.End()
	})

	Test(t, "replace-test-message: transform: wrong-fixture", func(t *T) {
		t.Transform("wrong-fixture")
		t.End()
	})
}
