//go:build ignore

package fixture

func f(a, b bool) {
	if a {
		println("a")
	} else if b {
		println("b")
	}
}
