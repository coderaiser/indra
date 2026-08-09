package conditions_test

import (
	"testing"

	. "coderaiser/indra/internal/plugin_conditions"

	. "github.com/coderaiser/go-tape"
)

func TestRules(t *testing.T) {
	Test(t, "plugin-conditions: remove-useless-comments is registered", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "remove-useless-comments" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})
}
