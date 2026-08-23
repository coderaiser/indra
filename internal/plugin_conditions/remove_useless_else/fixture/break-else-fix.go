//go:build ignore

package fixture

func f(items []int) {
	for _, x := range items {
		if x > 0 {
			break
		}
		println(x)
	}
}
