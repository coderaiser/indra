package remove_useless_match

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ ast.Node) string { return "remove useless Match" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessMatch}
}

func findUselessMatch(node ast.Node, push func(ast.Node)) {
	file := node.(*ast.File)
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Match" {
			continue
		}
		if isUselessMatch(fn) {
			push(file)
			break
		}
	}
}

// Fix removes useless Match() decls from file.Decls. node is *ast.File;
// options is unused.
func Fix(node ast.Node, _ map[string]any) {
	file := node.(*ast.File)
	kept := file.Decls[:0]
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name != nil && fn.Name.Name == "Match" && isUselessMatch(fn) {
			continue
		}
		kept = append(kept, decl)
	}
	file.Decls = kept
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(node ast.Node) string            { return Report(node) }
func (Plugin) Traverse() Traverser                    { return Traverse() }
func (Plugin) Fix(node ast.Node, opts map[string]any) { Fix(node, opts) }
