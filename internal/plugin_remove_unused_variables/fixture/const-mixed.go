//go:build ignore

package fixture

const (
	timeout = 30
	used    = 1
)

func f() {
	_ = used
}

type unusedType struct{}