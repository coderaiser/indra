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

	Test(t, "tape: transform convert-equal-to-not-ok", func(t *indratest.T) {
		t.Transform("convert-equal-to-not-ok")
		t.End()
	})

	Test(t, "tape: transform convert-ok-to-not-ok", func(t *indratest.T) {
		t.Transform("convert-ok-to-not-ok")
		t.End()
	})

	Test(t, "tape: transform remove-useless-condition", func(t *indratest.T) {
		t.Transform("remove-useless-condition")
		t.End()
	})

	Test(t, "tape: transform remove-useless-prefix", func(t *indratest.T) {
		t.Transform("remove-useless-prefix")
		t.End()
	})
}
