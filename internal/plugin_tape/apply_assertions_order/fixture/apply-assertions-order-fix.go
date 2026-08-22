//go:build ignore

package fixture

import (
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {

	fmt.Println("gap")
	t.Equal(1, 1)

	t.End()
}
