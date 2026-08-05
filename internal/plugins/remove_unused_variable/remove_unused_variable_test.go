package remove_unused_variable_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestRemoveUnusedVariable(t *testing.T) {
	Test(t, "remove-unused-variable: report unused var", func(t *indratest.T) {
		t.Report("unused-var", "remove unused variable: x")
		t.End()
	})

	Test(t, "remove-unused-variable: no report when all vars used", func(t *indratest.T) {
		t.NoReport("used-var")
		t.End()
	})

	Test(t, "remove-unused-variable: no report for blank in tuple assign", func(t *indratest.T) {
		t.NoReport("tuple-blank")
		t.End()
	})

	Test(t, "remove-unused-variable: fix removes unused var", func(t *indratest.T) {
		t.Transform("unused-var")
		t.End()
	})

	Test(t, "remove-unused-variable: fix blank the unused var in tuple", func(t *indratest.T) {
		t.Transform("tuple-unused")
		t.End()
	})

	Test(t, "remove-unused-variable: fix drops tuple with blank and unused var", func(t *indratest.T) {
		t.Transform("tuple-blank-unused")
		t.End()
	})
}
