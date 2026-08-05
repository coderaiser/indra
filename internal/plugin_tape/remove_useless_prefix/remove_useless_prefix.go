package remove_useless_prefix

import (
	"go/ast"
	"reflect"

	. "coderaiser/indra/types"
)

func Report(_ ast.Node) string { return "remove useless tape prefix" }

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessPrefix}
}

func findUselessPrefix(node ast.Node, push func(ast.Node)) {
	file := node.(*ast.File)
	alias, _ := findTapeImport(file)
	if alias == "" {
		return
	}
	if hasLocalCollision(file, alias) {
		return
	}
	push(node)
}

// Fix rewrites a named go-tape import to a dot import and drops the alias
// prefix from every selector use (tape.X → X).
func Fix(node ast.Node, _ map[string]any) {
	file := node.(*ast.File)
	alias, spec := findTapeImport(file)
	if alias == "" {
		return
	}
	if hasLocalCollision(file, alias) {
		return
	}
	spec.Name = &ast.Ident{Name: ".", NamePos: spec.Name.NamePos}
	replaceSelectors(reflect.ValueOf(file), alias)
}
