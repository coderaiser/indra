package remove_useless_condition_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/remove_useless_condition"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("remove-useless-condition", remove_useless_condition.Plugin{})

func TestRemoveUselessCondition(t *testing.T) {
	Test(t, "remove-useless-condition: report Ok(err != nil)", func(t *T) {
		t.Report("remove-useless-condition", "remove useless condition")
		t.End()
	})

	Test(t, "remove-useless-condition: transform Ok(err != nil) to Ok(err)", func(t *T) {
		t.Transform("remove-useless-condition")
		t.End()
	})

	Test(t, "remove-useless-condition: no report for Ok(err)", func(t *T) {
		t.NoReport("ok-direct")
		t.End()
	})

	Test(t, "remove-useless-condition: no transform for Ok(err)", func(t *T) {
		t.NoTransform("ok-direct")
		t.End()
	})

	Test(t, "remove-useless-condition: report NotOk(err == nil)", func(t *T) {
		t.Report("not-ok-nil", "remove useless condition")
		t.End()
	})

	Test(t, "remove-useless-condition: transform NotOk(err == nil) to NotOk(err)", func(t *T) {
		t.Transform("not-ok-nil")
		t.End()
	})

	Test(t, "remove-useless-condition: no report for NotOk(err)", func(t *T) {
		t.NoReport("not-ok-direct")
		t.End()
	})

	Test(t, "remove-useless-condition: no transform for NotOk(err)", func(t *T) {
		t.NoTransform("not-ok-direct")
		t.End()
	})
}
