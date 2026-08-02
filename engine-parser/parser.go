package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Parse parses Go source bytes into an AST file and token set.
// Comments are preserved.
func Parse(src []byte) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return file, fset, nil
}
