//go:build ignore

package fixture

import "github.com/lithammer/dedent"

func f() []byte {
	return []byte(dedent.Dedent(`
		[ignore]
		patterns = ["vendor/**"]
		`))
}
