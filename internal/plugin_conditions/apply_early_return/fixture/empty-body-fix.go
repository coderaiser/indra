//go:build ignore

package fixture

func f(x int) {
	if x > 0 {
		return
	}
	println("non-positive")
}
