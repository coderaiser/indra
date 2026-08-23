//go:build ignore

package fixture

func f(a bool, x *int) {
	if v, ok := x.(*int); ok {
		println(v)
	} else if a {
		println(v)
	}
}
