//go:build ignore

package fixture

func f(t *testing.T) {
	Test(t, "report something", func(t *T) {
		t.End()
	})
}
