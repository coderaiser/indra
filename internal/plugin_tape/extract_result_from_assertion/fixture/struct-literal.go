//go:build ignore

package fixture

import "testing"

type point struct {
	x int
	y int
}

func TestFoo(t *testing.T) {
	t.Equal(got, point{1, 2})
	t.DeepEqual(want, point{3, 4})
	t.End()
}
