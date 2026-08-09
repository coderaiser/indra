//go:build ignore

package fixture

type T struct{}

func (t T) helper() {}

func (t *T) pointerHelper() {}
