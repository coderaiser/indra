//go:build ignore

package fixture

func f(a, b bool) {
	if a {
		println("same")
	} else if b {
		println("same")
	}
}
