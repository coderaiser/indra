//go:build ignore

package fixture

func f(items []int) int {
	g := func(items []int) int {
		for _, x := range items {
			if x < 0 {
				continue
			}
			return x
		}
		return 0
	}
	return g(items)
}
