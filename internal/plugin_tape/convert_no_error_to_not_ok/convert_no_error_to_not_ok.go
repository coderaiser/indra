package convert_no_error_to_not_ok

import (
	"go/ast"

	. "coderaiser/indra/types"
)

const goTapePath = `"github.com/coderaiser/go-tape"`

func Report(_ ast.Node) string { return "convert NoError(err) to NotOk(err)" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findNoErrorCalls}
}

func findNoErrorCalls(node ast.Node, push func(ast.Node)) {
	file := node.(*ast.File)
	if !hasGoTapeImport(file) {
		return
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if sel.Sel != nil && sel.Sel.Name == "NoError" {
			found = true
			return false
		}
		return true
	})
	if found {
		push(file)
	}
}

func Fix(node ast.Node, _ map[string]any) {
	file := node.(*ast.File)
	if !hasGoTapeImport(file) {
		return
	}
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if sel.Sel != nil && sel.Sel.Name == "NoError" {
			sel.Sel = &ast.Ident{Name: "NotOk"}
		}
		return true
	})
}

func hasGoTapeImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path != nil && imp.Path.Value == goTapePath {
			return true
		}
	}
	return false
}
