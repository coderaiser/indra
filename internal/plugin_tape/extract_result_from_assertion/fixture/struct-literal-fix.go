//go:build ignore

package fixture

import "testing"

type point struct {
	x int
	y int
}

func TestFoo(t *testing.T) {
	expected := point{1, 2}
	t.Equal(got, expected)
	expected := point{3, 4}
	t.DeepEqual(want, expected)

	t.End()
}
