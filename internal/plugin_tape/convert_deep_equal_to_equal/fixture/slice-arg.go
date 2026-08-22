//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.DeepEqual(list, []string{"a"})
	t.End()
}
