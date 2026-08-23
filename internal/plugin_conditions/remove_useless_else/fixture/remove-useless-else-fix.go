//go:build ignore

package fixture

func f(x int) string {
	if x > 0 {
		return "positive"
	}
	return "non-positive"
}
