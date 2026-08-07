//go:build ignore

package fixture

func report(x string) string {
	switch x {
	case "a":
		return "A"
	case "b":
		return "B"
	}
	return "unknown"
}
