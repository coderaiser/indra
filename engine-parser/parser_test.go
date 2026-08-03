package parser_test

import (
	"testing"

	parser "coderaiser/indra/engine-parser"
	tape "github.com/coderaiser/go-tape"
)

func TestParse(t *testing.T) {
	tape.Test(t, "parse: valid source returns non-nil file", func(t *tape.T) {
		file, _, _ := parser.Parse([]byte("package p\nfunc f() {}\n"))
		t.Ok(file != nil)
		t.End()
	})

	tape.Test(t, "parse: valid source returns non-nil fset", func(t *tape.T) {
		_, fset, _ := parser.Parse([]byte("package p\nfunc f() {}\n"))
		t.Ok(fset != nil)
		t.End()
	})

	tape.Test(t, "parse: valid source returns nil error", func(t *tape.T) {
		_, _, err := parser.Parse([]byte("package p\nfunc f() {}\n"))
		t.Equal(err, nil)
		t.End()
	})

	tape.Test(t, "parse: invalid source returns non-nil error", func(t *tape.T) {
		_, _, err := parser.Parse([]byte("package p\nfunc (\n"))
		t.Ok(err != nil)
		t.End()
	})
}
