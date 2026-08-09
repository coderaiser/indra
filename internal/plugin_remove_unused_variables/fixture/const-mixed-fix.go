//go:build ignore

package fixture

const (
	used = 1
)

func F() {
	_ = used
}

type unusedType struct{}
