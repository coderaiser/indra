//go:build ignore

package fixture

func f(x int) {
	if x > 0 {
		println("positive")
	} else {
		println("non-positive")
	}
	println("done")
}
