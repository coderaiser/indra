//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	Suite(t, "allowed-receiver: something", func(t *Suite.T) {
		t.Equal(1, 1)
	})

}
