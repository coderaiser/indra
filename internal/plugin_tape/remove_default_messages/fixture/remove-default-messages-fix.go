//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.Ok(cond)
	t.NotOk(other)
	t.Equal(a, b)
	t.DeepEqual(c, d)

	t.End()
}
