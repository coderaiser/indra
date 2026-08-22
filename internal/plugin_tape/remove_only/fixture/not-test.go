//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	other.Only(t, "foo: something", func(t *other.T) {
		t.Equal(1, 1)
	})
}
