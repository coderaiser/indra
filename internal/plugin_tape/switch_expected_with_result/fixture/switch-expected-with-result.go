//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(expected, actual)
	t.End()
}
