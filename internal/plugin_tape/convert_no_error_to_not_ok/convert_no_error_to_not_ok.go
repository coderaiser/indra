// Package convert_no_error_to_not_ok rewrites tape's t.NoError(err) into the
// equivalent t.NotOk(err). Only files that import go-tape are considered, so the
// collision-prone NoError name is never rewritten outside tape.
package convert_no_error_to_not_ok

import (
	"go/ast"

	. "coderaiser/indra/types"
)

const goTapePath = `"github.com/coderaiser/go-tape"`

func Report(_ Path) string { return "convert NoError(err) to NotOk(err)" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findNoErrorCalls}
}

func findNoErrorCalls(p Path, push func(Path)) {
	file := p.Node.(*ast.File)
	if !hasGoTapeImport(file) {
		return
	}
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			call := callPath.Node.(*ast.CallExpr)
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			if sel.Sel != nil && sel.Sel.Name == "NoError" {
				found = true
				callPath.Stop()
			}
		},
	})
	if found {
		push(p)
	}
}

// Fix rewrites tape t.NoError(err) calls to t.NotOk(err). It is only ever
// invoked on a pushed *ast.File, so findNoErrorCalls has already guaranteed a
// go-tape import -- no import guard is needed here.
func Fix(p Path, _ map[string]any) {
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			call := callPath.Node.(*ast.CallExpr)
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			if sel.Sel != nil && sel.Sel.Name == "NoError" {
				sel.Sel = &ast.Ident{Name: "NotOk"}
			}
		},
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

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
