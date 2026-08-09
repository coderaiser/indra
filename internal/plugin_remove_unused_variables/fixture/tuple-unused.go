//go:build ignore

package fixture

func G() (int, int) { return 1, 2 }

func F() {
	a, b := G()
	_ = a
}
