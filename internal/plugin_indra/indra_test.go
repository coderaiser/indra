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

	Test(t, "plugin-indra: apply-fixture-name-to-message is present", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "apply-fixture-name-to-message" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "plugin-indra: replace-test-message is present", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "replace-test-message" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "plugin-indra: convert-inspect-to-traverse is present", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "convert-inspect-to-traverse" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "plugin-indra: apply-compare is present", func(t *T) {
		found := false
		for _, r := range Rules() {
			if r.Name == "apply-compare" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})
}
