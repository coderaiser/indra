//go:build ignore

package fixture

func f(items []int) int {
	for _, x := range items {
		if x < 0 {
			continue
		}
		return x
	}
	return 0
}
