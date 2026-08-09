//go:build ignore

package fixture

const (
	timeout = 30
	used    = 1
)

func F() {
	_ = used
}

type unusedType struct{}