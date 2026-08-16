package operator

import "go/ast"

// Path wraps an AST node with parent context.
type Path struct {
	Node ast.Node
	Parent ast.Node
	Field  string
	Index  int
}

// Remove removes path.Node from its parent.
func Remove(path *Path) {
	switch p := path.Parent.(type) {
	case *ast.BlockStmt:
		if path.Index >= 0 && path.Index < len(p.List) {
			p.List = append(p.List[:path.Index], p.List[path.Index+1:]...)
		}
	case *ast.GenDecl:
		if path.Index >= 0 && path.Index < len(p.Specs) {
			p.Specs = append(p.Specs[:path.Index], p.Specs[path.Index+1:]...)
		}
	case *ast.File:
		if path.Index >= 0 && path.Index < len(p.Decls) {
			p.Decls = append(p.Decls[:path.Index], p.Decls[path.Index+1:]...)
		}
	}
}

// remove is the method form of Remove.
func (path *Path) remove() {
	Remove(path)
}
