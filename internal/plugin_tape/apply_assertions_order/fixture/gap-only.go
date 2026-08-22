//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	x := 1
	_ = x
	t.End()
}
