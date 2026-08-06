package remove_unused_variables

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(node ast.Node) string {
	if node == nil {
		return "remove unused variable"
	}
	switch n := node.(type) {
	case *ast.File:
		imports := collectImports(n)
		used := countIdentUses(n)
		for _, imp := range imports {
			if imp.blank || imp.dot {
				continue
			}
			if used[imp.localName] == 0 {
				return "remove unused import: " + imp.path
			}
		}
		consts := unusedConstNames(n)
		if len(consts) > 0 {
			return "remove unused const: " + consts[0]
		}
	case *ast.BlockStmt:
		unused := unusedVarNames(n)
		if len(unused) > 0 {
			return "remove unused variable: " + unused[0]
		}
	}
	return "remove unused variable"
}

func Traverse() Traverser {
	return Traverser{
		"*ast.File":      findUnusedImportsAndConsts,
		"*ast.BlockStmt": findUnusedVars,
	}
}

func Fix(node ast.Node, _ map[string]any) {
	switch n := node.(type) {
	case *ast.File:
		fixUnusedImports(n)
		fixUnusedConsts(n)
	case *ast.BlockStmt:
		fixUnusedVars(n)
	}
}

type Plugin struct{}

func (Plugin) Report(node ast.Node) string            { return Report(node) }
func (Plugin) Traverse() Traverser                    { return Traverse() }
func (Plugin) Fix(node ast.Node, opts map[string]any) { Fix(node, opts) }
