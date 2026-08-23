//go:build ignore

package fixture

func f(x int) string {
	if x > 0 {
		println("positive")
	} else if x < 0 {
		println("negative")
	}
	return ""
}
