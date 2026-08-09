//go:build ignore

package fixture

func unusedHelper() string {
	return "hello"
}

func ExportedFunc() string {
	return "world"
}
