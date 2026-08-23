//go:build ignore

package fixture

func f(a, b bool, x *int) {
	if v, ok := x.(*int); ok {
		if b {
			println(v, a)
		}
	}
}
