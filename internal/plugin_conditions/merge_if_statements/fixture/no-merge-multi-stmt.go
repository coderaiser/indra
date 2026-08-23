//go:build ignore

package fixture

func f(a, b bool) {
	if a {
		println("a")
		if b {
			println("b")
		}
	}
}
