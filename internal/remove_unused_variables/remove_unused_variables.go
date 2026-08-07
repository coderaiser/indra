package remove_unused_variables

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(node ast.Node) string {
	switch n := node.(type) {
	case *importFinding:
		return "remove unused import: " + n.spec.Path.Value
	case *ast.File:
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
	case *importFinding:
		fixOneUnusedImport(n.file, n.spec)
	case *ast.File:
		fixUnusedConsts(n)
	case *ast.BlockStmt:
		fixUnusedVars(n)
	}
}

type Plugin struct{}

func (Plugin) Report(node ast.Node) string            { return Report(node) }
func (Plugin) Traverse() Traverser                    { return Traverse() }
func (Plugin) Fix(node ast.Node, opts map[string]any) { Fix(node, opts) }
