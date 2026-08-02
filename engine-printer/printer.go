package printer

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/token"
)

// Print formats an AST file back into Go source bytes.
// Returns error if file is nil or formatting fails.
func Print(file *ast.File, fset *token.FileSet) ([]byte, error) {
	if file == nil {
		return nil, errors.New("engine-printer: nil file")
	}
	if fset == nil {
		fset = token.NewFileSet()
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
