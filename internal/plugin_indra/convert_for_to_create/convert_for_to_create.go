package convert_for_to_create

import (
	"go/ast"

	. "coderaiser/indra/types"
)

const testImportPath = `"coderaiser/indra/internal/test"`

func Report(_ Path) string {
	return "convert indratest.For to CreateTest"
}

func Traverse() Traverser {
	return Traverser{"*ast.File": findForCalls}
}

func findForCalls(p Path, push func(Path)) {
	file, ok := p.Node.(*ast.File)
	if !ok {
		return
	}
	if hasForCall(file) {
		push(p)
	}
}

// hasForCall reports whether file references an indratest.For call.
func hasForCall(file *ast.File) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		id, ok := sel.X.(*ast.Ident)
		if ok && id.Name == "indratest" && sel.Sel != nil && sel.Sel.Name == "For" {
			found = true
			return false
		}
		return true
	})
	return found
}

// Fix rewrites indratest.For usages (call site, import alias, and *T type)
// to the createTest form. node is *ast.File; options is unused.
func Fix(p Path, _ map[string]any) {
	file, ok := p.Node.(*ast.File)
	if !ok {
		return
	}
	rewriteImport(file)
	rewriteCalls(file)
	rewriteT(file)
}

// rewriteImport turns `indratest "coderaiser/indra/internal/test"` into a dot
// import so the unqualified createTest name resolves.
func rewriteImport(file *ast.File) {
	for _, imp := range file.Imports {
		if imp.Path == nil || imp.Path.Value != testImportPath {
			continue
		}
		if imp.Name != nil && imp.Name.Name == "indratest" {
			imp.Name = ast.NewIdent(".")
		}
	}
}

// rewriteCalls replaces `indratest.For(__a, __b)` call expressions with
// `createTest(__a, __b)`.
func rewriteCalls(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		id, ok := sel.X.(*ast.Ident)
		if ok && id.Name == "indratest" && sel.Sel != nil && sel.Sel.Name == "For" {
			call.Fun = ast.NewIdent("CreateTest")
		}
		return true
	})
}

// rewriteT replaces `*indratest.T` callback parameters with `*T`.
func rewriteT(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		star, ok := n.(*ast.StarExpr)
		if !ok {
			return true
		}
		sel, ok := star.X.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		id, ok := sel.X.(*ast.Ident)
		if ok && id.Name == "indratest" && sel.Sel != nil && sel.Sel.Name == "T" {
			star.X = ast.NewIdent("T")
		}
		return true
	})
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
