//go:build ignore

package fixture

func f(x int) {
	g := func(x bool) {
		if x {
			println(x)
		}
	}
	g(x)
}
