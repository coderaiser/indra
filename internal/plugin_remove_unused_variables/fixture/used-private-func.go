//go:build ignore

package fixture

func helper() string {
	return "hello"
}

func ExportedFunc() string {
	return helper()
}
