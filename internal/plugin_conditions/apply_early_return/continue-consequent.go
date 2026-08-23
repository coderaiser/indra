//go:build ignore

package fixture

func f(items []int) int {
	total := 0
	for _, x := range items {
		if x < 0 {
			continue
		} else {
			total += x
		}
	}
	return total
}
