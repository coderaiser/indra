//go:build ignore

package fixture

func TestEmptyString(t *testing.T) {
	Test(t, "convert-equal-to-not-ok: transform: empty-string", func(t *T) {
		t.NotOk(out)

		t.End()
	})
}
