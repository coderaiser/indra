//go:build ignore

package fixture

import "github.com/coderaiser/go-tape"

type MyType int

func f(err error) {
	t.Equal(result, foo())
	t.Equal(result, true)
	t.Equal(result, undefined)
	t.Equal(result, err)
	t.Equal(result, MyType)
}
