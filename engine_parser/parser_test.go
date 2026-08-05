package parser_test

import (
	"testing"

	parser "coderaiser/indra/engine_parser"

	. "github.com/coderaiser/go-tape"
)

func TestParse(t *testing.T) {
	Test(t, "parse: valid source returns non-nil file", func(t *T) {
		file, _, _ := parser.Parse([]byte("package p\nfunc f() {}\n"))
		t.Ok(file)

		t.End()
	})

	Test(t, "parse: valid source returns non-nil fset", func(t *T) {
		_, fset, _ := parser.Parse([]byte("package p\nfunc f() {}\n"))
		t.Ok(fset)

		t.End()
	})

	Test(t, "parse: valid source returns nil error", func(t *T) {
		_, _, error := parser.Parse([]byte("package p\nfunc f() {}\n"))
		t.NotOk(error)
		t.End()
	})

	Test(t, "parse: invalid source returns non-nil error", func(t *T) {
		_, _, error := parser.Parse([]byte("package p\nfunc (\n"))
		t.Ok(error)
		t.End()
	})
}
