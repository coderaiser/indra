package main

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestRegistry(t *testing.T) {
	tape.Test(t, "registry: Registry is non-empty", func(t *tape.T) {
		t.Ok(len(Registry) > 0)
		t.End()
	})

	tape.Test(t, "registry: tape group carries rules", func(t *tape.T) {
		found := false
		for _, p := range Registry {
			if p.Name == "tape" {
				found = p.Rules != nil && len(p.Rules) > 0
			}
		}
		t.Ok(found)
		t.End()
	})

	tape.Test(t, "registry: every entry has a name", func(t *tape.T) {
		ok := true
		for _, p := range Registry {
			if p.Name == "" {
				ok = false
			}
		}
		t.Ok(ok)
		t.End()
	})

	tape.Test(t, "registry: every entry has either rules or a plugin", func(t *tape.T) {
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
