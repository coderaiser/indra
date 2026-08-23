//go:build ignore

package fixture

var g = func(x int) {
	if x > 0 {
		println("positive")
	} else {
		println("non-positive")
	}
}
