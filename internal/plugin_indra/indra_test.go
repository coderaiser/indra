package plugin_indra_test

import (
	"testing"

	. "coderaiser/indra/internal/plugin_indra"

	. "github.com/coderaiser/go-tape"
)

func TestRules(t *testing.T) {
	Test(t, "plugin-indra: remove-useless-match is present", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "remove-useless-match" {
				found = true
			}
		}
		t.Equal(found, true)
		t.End()
	})

	Test(t, "plugin-indra: remove-useless-match is Disabled by default", func(t *T) {
		for _, r := range Rules() {
			if r.Name == "remove-useless-match" {
				t.Equal(r.Disabled, true)
			}
		}
		t.End()
	})
}
