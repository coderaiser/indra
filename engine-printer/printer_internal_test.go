package printer

import (
	"errors"
	"go/ast"
	"go/token"
	"io"
	"testing"
)

// TestPrintFormatError covers the format.Node failure branch by temporarily
// swapping the hoisted formatNode var to always return an error.
func TestPrintFormatError(t *testing.T) {
	orig := formatNode
	formatNode = func(_ io.Writer, _ *token.FileSet, _ any) error {
		return errors.New("boom")
	}
	defer func() { formatNode = orig }()

	_, err := Print(&ast.File{Name: &ast.Ident{Name: "p"}}, token.NewFileSet())
	if err == nil {
		t.Fatal("expected formatting error, got nil")
	}
}
