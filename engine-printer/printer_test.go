package printer_test

import (
	"testing"

	engineparser  "coderaiser/indra/engine-parser"
	engineprinter "coderaiser/indra/engine-printer"
	tape "github.com/coderaiser/go-tape"
)

func TestPrint(t *testing.T) {
	tape.Test(t, "printer: roundtrip preserves source", func(t *tape.T) {
		src := []byte("package p\n\nfunc f() {}\n")
		file, fset, _ := engineparser.Parse(src)
		got, _ := engineprinter.Print(file, fset)
		t.Equal(string(got), string(src))
		t.End()
	})

	tape.Test(t, "printer: nil file returns error", func(t *tape.T) {
		_, error := engineprinter.Print(nil, nil)
		t.Ok(error)
		t.End()
	})

	tape.Test(t, "printer: nil fset still prints", func(t *tape.T) {
		src := []byte("package p\n\nfunc f() {}\n")
		file, _, _ := engineparser.Parse(src)
		got, _ := engineprinter.Print(file, nil)
		t.Ok(len(got) > 0)
		t.End()
	})
}

