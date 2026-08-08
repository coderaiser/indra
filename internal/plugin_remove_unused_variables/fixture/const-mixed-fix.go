//go:build ignore

package fixture

const (
	used = 1
)

func f() {
	_ = used
}

type unusedType struct{}
