//go:build ignore

package fixture

import (
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	t.Equal(1, 1)
	fmt.Println("gap")
	t.End()
}
