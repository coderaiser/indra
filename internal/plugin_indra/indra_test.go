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
		t.Ok(found)

		t.End()
	})

	Test(t, "plugin-indra: remove-useless-match is Disabled by default", func(t *T) {
		for _, r := range Rules() {
			if r.Name == "remove-useless-match" {
				t.Ok(r.Disabled)

			}
		}
		t.End()
	})

	Test(t, "plugin-indra: convert-for-to-create-test is present", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "convert-for-to-create-test" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "plugin-indra: convert-for-to-create-test is Disabled by default", func(t *T) {
		for _, r := range Rules() {
			if r.Name == "convert-for-to-create-test" {
				t.Ok(r.Disabled)
			}
		}
		t.End()
	})
}
