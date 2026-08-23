//go:build ignore

package fixture

func f() bool {
	type B bool
	var b B
	return b == true
}
