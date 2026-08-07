package remove_useless_prefix

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove useless tape prefix" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessPrefix}
}

func findUselessPrefix(p Path, push func(Path)) {
	file, ok := p.Node.(*ast.File)
	if !ok {
		return
	}
	alias, _ := findTapeImport(file)
	if alias == "" {
		return
	}
	if hasLocalCollision(file, alias) {
		return
	}
	push(p)
}

// Fix rewrites a named go-tape import to a dot import and drops the alias
// prefix from every selector use (tape.X → X).
func Fix(p Path, _ map[string]any) {
	file, ok := p.Node.(*ast.File)
	if !ok {
		return
	}
	alias, spec := findTapeImport(file)
	if alias == "" {
		return
	}
	if hasLocalCollision(file, alias) {
		return
	}
	spec.Name = &ast.Ident{Name: ".", NamePos: spec.Name.NamePos}

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		if sel, ok := c.Node().(*ast.SelectorExpr); ok && selMatchesAlias(sel, alias) {
			c.Replace(sel.Sel)
		}
		return true
	}, nil)
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
