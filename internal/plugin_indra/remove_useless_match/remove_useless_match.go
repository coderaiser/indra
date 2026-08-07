package remove_useless_match

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove useless Match" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessMatch}
}

func findUselessMatch(p Path, push func(Path)) {
	file, ok := p.Node.(*ast.File)
	if !ok {
		return
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Match" {
			continue
		}
		if isUselessMatch(fn) {
			push(p)
			break
		}
	}
}

// Fix removes useless Match() decls from file.Decls. node is *ast.File;
// options is unused.
func Fix(p Path, _ map[string]any) {
	file, ok := p.Node.(*ast.File)
	if !ok {
		return
	}
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

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
