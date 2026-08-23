package conditions_test

import (
	"testing"

	. "coderaiser/indra/internal/plugin_conditions"

	. "github.com/coderaiser/go-tape"
)

func TestRules(t *testing.T) {
	Test(t, "conditions: Rules returns 9 rules", func(t *T) {
		t.Equal(len(Rules()), 9)
		t.End()
	})

	Test(t, "conditions: remove-useless-comments is registered", func(t *T) {
		t.Equal(Rules()[0].Name, "remove-useless-comments")
		t.End()
	})

	Test(t, "conditions: convert-switch-to-if is registered", func(t *T) {
		t.Equal(Rules()[1].Name, "convert-switch-to-if")
		t.End()
	})

	Test(t, "conditions: remove-boolean is registered", func(t *T) {
		t.Equal(Rules()[2].Name, "remove-boolean")
		t.End()
	})

	Test(t, "conditions: reverse-condition is registered", func(t *T) {
		t.Equal(Rules()[3].Name, "reverse-condition")
		t.End()
	})

	Test(t, "conditions: remove-useless-else is registered", func(t *T) {
		t.Equal(Rules()[4].Name, "remove-useless-else")
		t.End()
	})

	Test(t, "conditions: merge-if-statements is registered", func(t *T) {
		t.Equal(Rules()[5].Name, "merge-if-statements")
		t.End()
	})

	Test(t, "conditions: merge-if-with-else is registered", func(t *T) {
		t.Equal(Rules()[6].Name, "merge-if-with-else")
		t.End()
	})

	Test(t, "conditions: apply-early-return is registered", func(t *T) {
		t.Equal(Rules()[7].Name, "apply-early-return")
		t.End()
	})

	Test(t, "conditions: simplify is registered", func(t *T) {
		t.Equal(Rules()[8].Name, "simplify")
		t.End()
	})
}
