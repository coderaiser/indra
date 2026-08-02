package printer_test

import (
	"testing"
	engineparser  "coderaiser/indra/engine-parser"
	engineprinter "coderaiser/indra/engine-printer"
)

func TestPrintRoundtrip(t *testing.T) {
	src := []byte("package p\n\nfunc f() {}\n")
	file, fset, err := engineparser.Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got, err := engineprinter.Print(file, fset)
	if err != nil {
		t.Fatalf("print: %v", err)
	}
	if string(got) != string(src) {
		t.Fatalf("roundtrip mismatch:\ngot:  %q\nwant: %q", got, src)
	}
}

func TestPrintNilFile(t *testing.T) {
	_, err := engineprinter.Print(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil file, got nil")
	}
}

func TestPrintNilFset(t *testing.T) {
	src := []byte("package p\n\nfunc f() {}\n")
	file, _, err := engineparser.Parse(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got, err := engineprinter.Print(file, nil)
	if err != nil {
		t.Fatalf("print: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected non-empty output")
	}
}

