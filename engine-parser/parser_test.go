package parser_test

import (
	"testing"
	parser "coderaiser/indra/engine-parser"
)

func TestParseValid(t *testing.T) {
	src := []byte("package p\nfunc f() {}\n")
	file, fset, err := parser.Parse(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file == nil {
		t.Fatal("expected non-nil file")
	}
	if fset == nil {
		t.Fatal("expected non-nil fset")
	}
}

func TestParseInvalid(t *testing.T) {
	_, _, err := parser.Parse([]byte("package p\nfunc (\n"))
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}
