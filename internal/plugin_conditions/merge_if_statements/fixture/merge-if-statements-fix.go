//go:build ignore

package fixture

func f(a, b bool) {
	if a && b {
		println("both")
	}
}
