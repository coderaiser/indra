//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	other.Skip(t, "foo: something", func(t *T) {
		t.Equal(1, 1)
	})
}
