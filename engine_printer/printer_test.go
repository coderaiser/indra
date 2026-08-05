package printer_test

import (
	"testing"

	engineparser "coderaiser/indra/engine_parser"
	engineprinter "coderaiser/indra/engine_printer"

	. "github.com/coderaiser/go-tape"
)

func TestPrint(t *testing.T) {
	Test(t, "printer: roundtrip preserves source", func(t *T) {
		src := []byte("package p\n\nfunc f() {}\n")
		file, fset, _ := engineparser.Parse(src)
		got, _ := engineprinter.Print(file, fset)
		result := string(got)
		t.Equal(result, string(src))

		t.End()
	})

	Test(t, "printer: nil file returns error", func(t *T) {
		_, error := engineprinter.Print(nil, nil)
		t.Ok(error)
		t.End()
	})

	Test(t, "printer: nil fset still prints", func(t *T) {
		src := []byte("package p\n\nfunc f() {}\n")
		file, _, _ := engineparser.Parse(src)
		got, _ := engineprinter.Print(file, nil)
		t.Ok(len(got) > 0)
		t.End()
	})
}
