package apply_fixture_name_to_message_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/apply_fixture_name_to_message"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("apply-fixture-name-to-message", apply_fixture_name_to_message.Plugin{})

func TestApplyFixtureNameToMessage(t *testing.T) {
	Test(t, "apply-fixture-name-to-message: report apply-fixture-name-to-message", func(t *T) {
		t.Report("apply-fixture-name-to-message", "apply fixture name to message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform apply-fixture-name-to-message", func(t *T) {
		t.Transform("apply-fixture-name-to-message")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no report already-prefixed", func(t *T) {
		t.NoReport("already-prefixed")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no report no-create-test", func(t *T) {
		t.NoReport("no-create-test")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no report odd-specs", func(t *T) {
		t.NoReport("odd-specs")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: no transform odd-specs", func(t *T) {
		t.NoTransform("odd-specs")
		t.End()
	})

	Test(t, "apply-fixture-name-to-message: transform mixed-message", func(t *T) {
		t.Transform("mixed-message")
		t.End()
	})
}
