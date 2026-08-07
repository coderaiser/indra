package remove_unused_variables

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(p Path) string {
	switch n := p.Node.(type) {
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

func Fix(p Path, _ map[string]any) {
	switch n := p.Node.(type) {
	case *importFinding:
		fixOneUnusedImport(n.file, n.spec)
	case *ast.File:
		fixUnusedConsts(n)
	case *ast.BlockStmt:
		fixUnusedVars(n)
	}
}

type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
