package plugins_test

import (
	"testing"

	"coderaiser/indra/internal/plugins"

	. "github.com/coderaiser/go-tape"
)

func TestRegistry(t *testing.T) {
	Test(t, "plugins: Registry is non-empty", func(t *T) {
		t.Ok(len(plugins.Registry) > 0)
		t.End()
	})

	Test(t, "plugins: tape group carries rules", func(t *T) {
		found := false
		for _, p := range plugins.Registry {
			if p.Name == "tape" {
				found = p.Rules != nil && len(p.Rules) > 0
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "plugins: every entry has a name", func(t *T) {
		found := true
		for _, p := range plugins.Registry {
			if p.Name == "" {
				found = false
			}
		}
		t.Ok(found)
		t.End()
	})

	Test(t, "plugins: every entry has either rules or a plugin", func(t *T) {
		ok := true
		for _, p := range plugins.Registry {
			if p.Rules == nil && p.Plugin == nil {
				ok = false
			}
		}
		t.Ok(ok)
		t.End()
	})
}
