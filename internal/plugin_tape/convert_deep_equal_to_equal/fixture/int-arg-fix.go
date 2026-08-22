//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(count, 42)

	t.End()
}
