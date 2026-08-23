//go:build ignore

package fixture

func f(x int) string {
	if x > 0 {
		return "positive"
	} else if x < 0 {
		return "negative"
	}
	return ""
}
