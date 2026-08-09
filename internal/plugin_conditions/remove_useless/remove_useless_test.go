package remove_useless_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/remove_useless"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/remove-useless", remove_useless.Plugin{})

func TestRemoveUselessCondition(t *testing.T) {
	Test(t, "remove-useless: report Ok(err != nil) conditions: remove-useless", func(t *T) {
		t.Report("remove-useless", "Avoid useless conditions")
		t.End()
	})

	Test(t, "conditions: remove-useless: transform Ok(err != nil) to Ok(err) conditions: remove-useless", func(t *T) {
		t.Transform("conditions: remove-useless")
		t.End()
	})

	Test(t, "conditions: remove-useless: no report for Ok(err) ok-direct", func(t *T) {
		t.NoReport("ok-direct")
		t.End()
	})

	Test(t, "conditions: remove-useless: no transform for Ok(err) ok-direct", func(t *T) {
		t.NoTransform("ok-direct")
		t.End()
	})

	Test(t, "conditions: remove-useless: report NotOk(err == nil) not-ok-nil", func(t *T) {
		t.Report("not-ok-nil", "remove useless condition")
		t.End()
	})

	Test(t, "conditions: remove-useless: transform NotOk(err == nil) to NotOk(err) not-ok-nil", func(t *T) {
		t.Transform("not-ok-nil")
		t.End()
	})

	Test(t, "conditions: remove-useless: no report for NotOk(err) not-ok-direct", func(t *T) {
		t.NoReport("not-ok-direct")
		t.End()
	})

	Test(t, "conditions: remove-useless: no transform for NotOk(err) not-ok-direct", func(t *T) {
		t.NoTransform("not-ok-direct")
		t.End()
	})

	Test(t, "conditions: remove-useless: no transform for NotOk(err) not-ok-direct", func(t *T) {
		t.Transform("false")
		t.End()
	})
}
