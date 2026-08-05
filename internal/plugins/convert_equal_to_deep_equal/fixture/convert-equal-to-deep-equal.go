//go:build ignore

package fixture

// convert-equal-to-deep-equal is the canonical happy path: Equal used on a slice.
func f() {
	t.Equal(x, []Block{})
}
