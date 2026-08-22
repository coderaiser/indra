//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	setup()
	t.End()
}
