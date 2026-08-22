//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.Ok(cond, "my custom message")
	t.Equal(a, b, "another custom one")
	t.End()
}
