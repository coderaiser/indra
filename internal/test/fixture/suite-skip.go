//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	Suite.Skip(t, "suite: something", func(t *Suite.T) {
		t.End()
	})
}
