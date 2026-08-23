//go:build ignore

package fixture

func f() bool {
	ok, no := true, false
	return no != true
}
