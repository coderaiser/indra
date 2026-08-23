//go:build ignore

package fixture

func f(x int) {
	g := func(x int) {
		if x > 0 {
			println("positive")
			return
		}
		println("non-positive")
	}
	g(x)
}
