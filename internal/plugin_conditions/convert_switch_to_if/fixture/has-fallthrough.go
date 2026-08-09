//go:build ignore

package fixture

func f(x string) string {
	switch x {
	case "a":
		fallthrough
	case "b":
		return "AB"
	}
	return ""
}
