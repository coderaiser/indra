package switch_expected_with_result_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/switch_expected_with_result"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("switch-expected-with-result", switch_expected_with_result.Plugin{})

func TestSwitchExpectedWithResult(t *testing.T) {
	Test(t, "switch-expected-with-result: report: switch-expected-with-result", func(t *T) {
		t.Report("switch-expected-with-result", `"result" should be before "expected"`)
		t.End()
	})

	Test(t, "switch-expected-with-result: transform: switch-expected-with-result", func(t *T) {
		t.Transform("switch-expected-with-result")
		t.End()
	})

	Test(t, "switch-expected-with-result: no report: already-correct", func(t *T) {
		t.NoReport("already-correct")
		t.End()
	})

	Test(t, "switch-expected-with-result: no transform: already-correct", func(t *T) {
		t.NoTransform("already-correct")
		t.End()
	})

	Test(t, "switch-expected-with-result: report: with-message", func(t *T) {
		t.Report("with-message", `"result" should be before "expected"`)
		t.End()
	})

	Test(t, "switch-expected-with-result: transform: with-message", func(t *T) {
		t.Transform("with-message")
		t.End()
	})

	Test(t, "switch-expected-with-result: report: deep-equal", func(t *T) {
		t.Report("deep-equal", `"result" should be before "expected"`)
		t.End()
	})

	Test(t, "switch-expected-with-result: transform: deep-equal", func(t *T) {
		t.Transform("deep-equal")
		t.End()
	})
}
