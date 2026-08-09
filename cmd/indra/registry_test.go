package main

import (
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestRegistry(t *testing.T) {
	Test(t, "registry: Registry is non-empty", func(t *T) {
		t.Ok(len(Registry) > 0)
		t.End()
	})

	Test(t, "registry: tape group carries rules", func(t *T) {
		found := false
		for _, p := range Registry {
			if p.Name == "tape" {
				found = len(p.Rules) > 0
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "registry: conditions group is registered", func(t *T) {
		found := false
		for _, p := range Registry {
			if p.Name == "conditions" {
				found = true
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "registry: every entry has a name", func(t *T) {
		ok := true
		for _, p := range Registry {
			if p.Name == "" {
				ok = false
			}
		}
		t.Ok(ok)
		t.End()
	})

	Test(t, "registry: every entry has either rules or a plugin", func(t *T) {
		ok := true
		for _, p := range Registry {
			if p.Rules == nil && p.Plugin == nil {
				ok = false
			}
		}
		t.Ok(ok)
		t.End()
	})
}
