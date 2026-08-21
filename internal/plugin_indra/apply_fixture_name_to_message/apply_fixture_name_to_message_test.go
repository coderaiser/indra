package apply_fixture_name_to_message_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/apply_fixture_name_to_message"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("apply-fixture-name-to-message", apply_fixture_name_to_message.Plugin{})

func TestApplyFixtureNameToMessage(t *testing.T) {
	Test(t, "apply-fixture-name-to-message: report: apply-fixture-name-to-message", func(t *T) {
		t.Report("apply-fixture-name-to-message", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: apply-fixture-name-to-message", func(t *T) {
		t.Transform("apply-fixture-name-to-message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no report: already-prefixed", func(t *T) {
		t.NoReport("already-prefixed")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no report: no-create-test", func(t *T) {
		t.NoReport("no-create-test")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no report: odd-specs", func(t *T) {
		t.NoReport("odd-specs")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no transform: odd-specs", func(t *T) {
		t.NoTransform("odd-specs")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: mixed-message", func(t *T) {
		t.Transform("mixed-message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: report: short-calls", func(t *T) {
		t.Report("short-calls", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: short-calls", func(t *T) {
		t.Transform("short-calls")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: report: prefix-only", func(t *T) {
		t.Report("prefix-only", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: prefix-only", func(t *T) {
		t.Transform("prefix-only")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: report: both-missing", func(t *T) {
		t.Report("both-missing", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: both-missing", func(t *T) {
		t.Transform("both-missing")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: report: direct-call", func(t *T) {
		t.Report("direct-call", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: direct-call", func(t *T) {
		t.Transform("direct-call")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: report: non-literal-report", func(t *T) {
		t.Report("non-literal-report", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: non-literal-report", func(t *T) {
		t.Transform("non-literal-report")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: report: non-func-callback", func(t *T) {
		t.Report("non-func-callback", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform: non-func-callback", func(t *T) {
		t.Transform("non-func-callback")
		t.End()
	})
}
