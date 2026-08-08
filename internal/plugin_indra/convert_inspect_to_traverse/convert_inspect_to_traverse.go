package convert_inspect_to_traverse

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "convert ast.Inspect to path.Traverse" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findInspectCalls}
}

// findInspectCalls pushes a Path for every ast.Inspect call expression inside a
// file that declares a Traverse() Traverser — i.e. an AST-walking plugin. Files
// without a Traverse marker (engine internals like engine_runner) are skipped.
func findInspectCalls(p Path, push func(Path)) {
	file := p.Node.(*ast.File)
	if !hasTraverser(file) {
		return
	}
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			if isInspectCall(callPath.Node.(*ast.CallExpr)) {
				push(callPath)
			}
		},
	})
}

// isInspectCall reports whether call is an ast.Inspect(f, fn) call expression.
func isInspectCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || id.Name != "ast" {
		return false
	}
	return sel.Sel != nil && sel.Sel.Name == "Inspect"
}

// hasTraverser reports whether the file declares a function returning the
// Traverser type — the marker of an AST-walking plugin file.
func hasTraverser(file *ast.File) bool {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if returnsTraverser(fn.Type) {
			return true
		}
	}
	return false
}

func returnsTraverser(ft *ast.FuncType) bool {
	if ft.Results == nil {
		return false
	}
	for _, r := range ft.Results.List {
		if id, ok := r.Type.(*ast.Ident); ok && id.Name == "Traverser" {
			return true
		}
	}
	return false
}

// Fix is a no-op: converting ast.Inspect to path.Traverse is report-only.
func Fix(Path, map[string]any) {
	_ = "report only"
}

// Plugin wraps the rule for the registry: a report-only AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
