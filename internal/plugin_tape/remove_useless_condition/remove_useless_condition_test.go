package remove_useless_condition_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_useless_condition"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("tape/remove-useless-condition", remove_useless_condition.Plugin{})

func TestRemoveUselessCondition(t *testing.T) {
	Test(t, "tape/remove-useless-condition: report: ok-nil", func(t *T) {
		t.Report("ok-nil", "remove useless condition")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: transform: ok-nil", func(t *T) {
		t.Transform("ok-nil")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: report: not-ok-nil", func(t *T) {
		t.Report("not-ok-nil", "remove useless condition")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: transform: not-ok-nil", func(t *T) {
		t.Transform("not-ok-nil")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: report: ok-false", func(t *T) {
		t.Report("ok-false", "remove useless condition")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: transform: ok-false", func(t *T) {
		t.Transform("ok-false")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: no report: no-condition", func(t *T) {
		t.NoReport("no-condition")
		t.End()
	})
}

// TestRemoveUselessConditionMessages covers the 3-arg forms where an assertion
// carries a message string: Ok(err != nil, msg) -> Ok(err, msg) and the
// NotOk == nil equivalent.
func TestRemoveUselessConditionMessages(t *testing.T) {
	Test(t, "tape/remove-useless-condition: report: ok-message", func(t *T) {
		t.Report("ok-message", "remove useless condition")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: transform: ok-message", func(t *T) {
		t.Transform("ok-message")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: report: not-ok-message", func(t *T) {
		t.Report("not-ok-message", "remove useless condition")
		t.End()
	})

	Test(t, "tape/remove-useless-condition: transform: not-ok-message", func(t *T) {
		t.Transform("not-ok-message")
		t.End()
	})
}
