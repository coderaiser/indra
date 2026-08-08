//go:build ignore

package fixture

func g(x int) int {
	if x == 1 {
		if x == 2 {
			return 2
		}
	}
	return 0
}
