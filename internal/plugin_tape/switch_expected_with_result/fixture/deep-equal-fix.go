//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.DeepEqual(actual, expected)
	t.Equal(actual, expected, "msg")

	t.End()
}
