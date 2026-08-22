//go:build ignore

package fixture

import "testing"

var counter int

func TestFoo(t *testing.T) {
	counter
	t.End()
}
