package remove_useless_condition_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestRemoveUselessCondition(t *testing.T) {
	Test(t, "remove-useless-condition: report Ok(err != nil)", func(t *indratest.T) {
		t.Report("remove-useless-condition", "remove useless condition: Ok(err != nil) → Ok(err)")
		t.End()
	})

	Test(t, "remove-useless-condition: transform Ok(err != nil) to Ok(err)", func(t *indratest.T) {
		t.Transform("remove-useless-condition")
		t.End()
	})

	Test(t, "remove-useless-condition: no report for Ok(err)", func(t *indratest.T) {
		t.NoReport("ok-direct")
		t.End()
	})

	Test(t, "remove-useless-condition: no transform for Ok(err)", func(t *indratest.T) {
		t.NoTransform("ok-direct")
		t.End()
	})
}
