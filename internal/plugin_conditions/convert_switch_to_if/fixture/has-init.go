//go:build ignore

package fixture

func f() string {
	switch x := "a"; x {
	case "a":
		return "A"
	}
	return ""
}
