package plugins_test

import (
	"testing"

	"coderaiser/indra/internal/plugins"

	. "github.com/coderaiser/go-tape"
)

func TestRegistries(t *testing.T) {
	Test(t, "plugins: All is non-empty", func(t *T) {
		t.Ok(len(plugins.All) > 0)
		t.End()
	})

	Test(t, "plugins: LoadInput contains All and Providers", func(t *T) {
		input := plugins.LoadInput()
		t.Ok(len(input) >= len(plugins.All)+len(plugins.Providers))
		t.End()
	})

	Test(t, "plugins: LoadInput is a fresh slice (does not alias All)", func(t *T) {
		origLen := len(plugins.All)
		input := plugins.LoadInput()
		// Appending to a returned slice must not mutate All.
		input = append(input, plugins.All[0])
		result := len(plugins.All)
		t.Equal(result, origLen)

		t.End()
	})

	Test(t, "plugins: every Provider path is unique", func(t *T) {
		seen := map[string]bool{}
		dup := false
		for _, p := range plugins.Providers {
			if seen[p.Path] {
				dup = true
			}
			seen[p.Path] = true
		}
		t.Ok(!dup)
		t.End()
	})
}
