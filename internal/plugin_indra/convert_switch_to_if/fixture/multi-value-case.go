//go:build ignore

package fixture

func f(x string) string {
	switch x {
	case "a", "b":
		return "AB"
	}
	return ""
}
