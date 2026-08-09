//go:build ignore

package fixture

func report(x string) string {
	if x == "a" {
		return "A"
	}
	if x == "b" {
		return "B"
	}
	return "unknown"
}
