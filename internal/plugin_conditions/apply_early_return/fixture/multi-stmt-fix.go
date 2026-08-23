//go:build ignore

package fixture

func f(x int) {
	if x > 0 {
		println("positive")
		return
	}
	println("non-positive")
	println("again")
}
