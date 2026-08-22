//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	expected := []string{"a"}
	result := []string{"b"}
	t.Equal(result, expected)
	t.End()
}
