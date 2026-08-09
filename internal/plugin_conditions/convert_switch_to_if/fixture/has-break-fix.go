//go:build ignore

package fixture

func f(x string) string {
	if x == "a" {
		return "A"
	}
	return "unknown"
}
