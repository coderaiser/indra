//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.DeepEqual(expected, actual)
	t.Equal(expected, actual, "msg")
	t.End()
}
