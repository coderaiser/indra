package reverse_condition_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_conditions/reverse_condition"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("conditions/reverse-condition", reverse_condition.Plugin{})

func TestReverseCondition(t *testing.T) {
	Test(t, "conditions/reverse-condition: report: report: reverse-condition", func(t *T) {
		t.Report("reverse-condition", "reverse condition")
		t.End()
	})

	Test(t, "conditions/reverse-condition: transform: transform: reverse-condition", func(t *T) {
		t.Transform("reverse-condition")
		t.End()
	})

	Test(t, "conditions/reverse-condition: transform: transform: transform: de-morgan", func(t *T) {
		t.Transform("de-morgan")
		t.End()
	})

	Test(t, "conditions/reverse-condition: transform: transform: transform: de-morgan-and", func(t *T) {
		t.Transform("de-morgan-and")
		t.End()
	})

	Test(t, "conditions/reverse-condition: transform: transform: transform: less", func(t *T) {
		t.Transform("less")
		t.End()
	})

	Test(t, "conditions/reverse-condition: transform: transform: transform: greater-or-equal", func(t *T) {
		t.Transform("greater-or-equal")
		t.End()
	})

	Test(t, "conditions/reverse-condition: transform: transform: transform: less-or-equal", func(t *T) {
		t.Transform("less-or-equal")
		t.End()
	})

	Test(t, "conditions/reverse-condition: no report: no report: no report: no-reverse", func(t *T) {
		t.NoReport("no-reverse")
		t.End()
	})

	Test(t, "conditions/reverse-condition: no report: no report: no report: unary-minus", func(t *T) {
		t.NoReport("unary-minus")
		t.End()
	})

	Test(t, "conditions/reverse-condition: no report: no report: no report: not-ident", func(t *T) {
		t.NoReport("not-ident")
		t.End()
	})

	Test(t, "conditions/reverse-condition: no report: no report: no report: equality", func(t *T) {
		t.NoReport("equality")
		t.End()
	})
}
