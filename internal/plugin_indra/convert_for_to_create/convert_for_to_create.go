package convert_for_to_create

import (
	"go/ast"

	. "coderaiser/indra/types"
)

const testImportPath = `"coderaiser/indra/internal/test"`

func Report(_ Path) string {
	return "convert indratest.For to CreateTest"
}

// Fix rewrites indratest.For usages (call site, import alias, and *T type)
// to the createTest form. node is *ast.File; options is unused.
func Fix(p Path, _ map[string]any) {
	file := p.Node.(*ast.File)
	rewriteImport(file)
	rewriteCalls(p)
	rewriteT(p)
}

func Traverse() Traverser {
	return Traverser{"*ast.File": findForCalls}
}

func findForCalls(p Path, push func(Path)) {
	if hasForCall(p) {
		push(p)
	}
}

// hasForCall reports whether the file rooted at p references an indratest.For
// call.
func hasForCall(p Path) bool {
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.SelectorExpr": func(sp Path) {
			sel := sp.Node.(*ast.SelectorExpr)
			id, ok := sel.X.(*ast.Ident)
			if ok && id.Name == "indratest" && sel.Sel != nil && sel.Sel.Name == "For" {
				found = true
				sp.Stop()
			}
		},
	})
	return found
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
func rewriteCalls(p Path) {
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(cp Path) {
			call := cp.Node.(*ast.CallExpr)
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			id, ok := sel.X.(*ast.Ident)
			if ok && id.Name == "indratest" && sel.Sel != nil && sel.Sel.Name == "For" {
				call.Fun = ast.NewIdent("CreateTest")
			}
		},
	})
}

// rewriteT replaces `*indratest.T` callback parameters with `*T`.
func rewriteT(p Path) {
	p.Traverse(map[string]func(Path){
		"*ast.StarExpr": func(sp Path) {
			star := sp.Node.(*ast.StarExpr)
			sel, ok := star.X.(*ast.SelectorExpr)
			if !ok {
				return
			}
			id, ok := sel.X.(*ast.Ident)
			if ok && id.Name == "indratest" && sel.Sel != nil && sel.Sel.Name == "T" {
				star.X = ast.NewIdent("T")
			}
		},
	})
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
