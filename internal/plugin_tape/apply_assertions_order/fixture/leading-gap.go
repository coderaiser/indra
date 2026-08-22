//go:build ignore

package fixture

import (
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	fmt.Println("setup")
	t.Equal(1, 1)
	t.End()
}
