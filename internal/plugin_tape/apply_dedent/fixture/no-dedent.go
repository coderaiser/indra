//go:build ignore

package fixture

func f() []byte {
	return []byte("[ignore]\\npatterns = [\\\"vendor/**\\\"]\\n")
}
