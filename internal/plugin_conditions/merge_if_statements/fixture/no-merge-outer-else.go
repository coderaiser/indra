//go:build ignore

package fixture

func f(a, b bool) {
	if a {
		if b {
			println("both")
		}
	} else {
		println("not a")
	}
}
