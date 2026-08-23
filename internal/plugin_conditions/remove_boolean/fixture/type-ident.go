//go:build ignore

package fixture

type status struct{}

func f() bool {
	ok := status == true
	return ok
}
