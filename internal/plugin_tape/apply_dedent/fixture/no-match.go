//go:build ignore

package fixture

func f() {
	println("x")
	foo().Dedent("a")
	x.Dedent("b")
	dedent.Strip("c")
	dedent.Dedent(a, b)
	dedent.Dedent(s)
}
