package tape_test

import (
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var Test = indratest.CreateTest(runtime.Caller(0))

func TestTape(t *testing.T) {
	Test(t, "tape: transform remove-skip", func(t *indratest.T) {
		t.Transform("remove-skip")
		t.End()
	})

	Test(t, "tape: transform add-t-end", func(t *indratest.T) {
		t.Transform("add-t-end")
		t.End()
	})
}
