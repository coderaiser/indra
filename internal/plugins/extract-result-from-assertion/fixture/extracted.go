//go:build ignore

package fixture

func f() {
	result := someFunc(a, b)
	t.DeepEqual(result, expected)
}
