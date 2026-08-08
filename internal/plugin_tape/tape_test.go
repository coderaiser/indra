package tape_test

import (
	"testing"

	tape "coderaiser/indra/internal/plugin_tape"
	. "coderaiser/indra/internal/test"
)

var Test = ForGroup("tape", tape.Rules())

func TestTape(t *testing.T) {
	Test(t, "tape: transform remove-skip", func(t *T) {
		t.Transform("remove-skip")
		t.End()
	})

	Test(t, "tape: transform add-t-end", func(t *T) {
		t.Transform("add-t-end")
		t.End()
	})

	Test(t, "tape: transform convert-equal-to-not-ok", func(t *T) {
		t.Transform("convert-equal-to-not-ok")
		t.End()
	})

	Test(t, "tape: transform convert-ok-to-not-ok", func(t *T) {
		t.Transform("convert-ok-to-not-ok")
		t.End()
	})

	Test(t, "tape: transform remove-useless-prefix", func(t *T) {
		t.Transform("remove-useless-prefix")
		t.End()
	})
}
