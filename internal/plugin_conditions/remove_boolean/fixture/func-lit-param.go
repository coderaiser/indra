//go:build ignore

package fixture

func f(x int) {
	g := func(x bool) {
		if x == true {
			println(x)
		}
	}
	g(x)
}
