//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	t.DeepEqual(name, "hello", "my message")
	t.End()
}
