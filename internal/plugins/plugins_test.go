package plugins_test

import (
	"testing"

	"coderaiser/indra/internal/plugins"
	tape "github.com/coderaiser/go-tape"
)

func TestRegistries(t *testing.T) {
	tape.Test(t, "plugins: All is non-empty", func(t *tape.T) {
		t.Ok(len(plugins.All) > 0)
		t.End()
	})

	tape.Test(t, "plugins: LoadInput contains All and Providers", func(t *tape.T) {
		input := plugins.LoadInput()
		t.Ok(len(input) >= len(plugins.All)+len(plugins.Providers))
		t.End()
	})

	tape.Test(t, "plugins: LoadInput is a fresh slice (does not alias All)", func(t *tape.T) {
		origLen := len(plugins.All)
		input := plugins.LoadInput()
		// Appending to a returned slice must not mutate All.
		input = append(input, plugins.All[0])
		t.Equal(len(plugins.All), origLen)
		t.End()
	})

	tape.Test(t, "plugins: every Provider path is unique", func(t *tape.T) {
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
