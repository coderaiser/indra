//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(1, 1)
	t.Equal(2, 2)
	t.End()
}
