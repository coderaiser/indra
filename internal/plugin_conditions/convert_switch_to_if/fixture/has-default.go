//go:build ignore

package fixture

func f(x string) string {
	switch x {
	case "a":
		return "A"
	default:
		return "other"
	}
	return ""
}
