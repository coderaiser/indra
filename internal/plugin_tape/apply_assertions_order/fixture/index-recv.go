//go:build ignore

package fixture

import "testing"

func TestFoo(t *testing.T) {
	suites[0].Equal(1, 1)
	suites[0].End()
}
