//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.Ok(cond, "should be truthy")
	t.NotOk(other, "should be falsy")
	t.Equal(a, b, "should equal")
	t.DeepEqual(c, d, "should deep equal")
	t.End()
}
