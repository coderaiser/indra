// Package tapeguard reports whether a source file imports go-tape.
//
// It is the operator-backed replacement for the removed operator.HasImport:
// tape-replacer plugins declare a guard `tapeImported` that delegates here, so
// a pattern only fires inside a file that actually uses go-tape.
package tapeguard

import (
	"go/ast"

	"coderaiser/indra/types"
)

// goTapePath is the import path of the go-tape package.
const goTapePath = `"github.com/coderaiser/go-tape"`

// Imported reports whether the file containing path imports go-tape. It is
// shaped as a types.MatchFn so it can be used directly as a Replacer guard.
// It walks the ancestor stack to the enclosing *ast.File and checks its
// imports by path (not by alias), so it works for aliased, dot, blank and
// unaliased imports alike.
func Imported(_ types.Vars, path types.Path) bool {
	for i := len(path.Stack) - 1; i >= 0; i-- {
		if file, ok := path.Stack[i].(*ast.File); ok {
			return hasImport(file)
		}
	}
	return false
}

// hasImport reports whether file imports go-tape by path.
func hasImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path.Value == goTapePath {
			return true
		}
	}
	return false
}
