//go:build ignore

package fixture

func f(items []int) {
	for _, x := range items {
		if x > 0 {
			break
		} else {
			println(x)
		}
	}
}
