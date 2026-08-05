package remove_skip_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestRemoveSkip(t *testing.T) {
	Test(t, "remove-skip: report Test.Skip call", func(t *indratest.T) {
		t.Report("skip", "remove Test.Skip call")
		t.End()
	})

	Test(t, "remove-skip: transform Test.Skip to Test", func(t *indratest.T) {
		t.Transform("skip")
		t.End()
	})

	Test(t, "remove-skip: no report for Test call", func(t *indratest.T) {
		t.NoReport("no-skip")
		t.End()
	})

	Test(t, "remove-skip: no transform for Test call", func(t *indratest.T) {
		t.NoTransform("no-skip")
		t.End()
	})
}
