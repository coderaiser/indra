//go:build ignore

package fixture

func g(x int) int {
	switch x {
	case 1:
		switch x {
		case 2:
			return 2
		}
	}
	return 0
}
