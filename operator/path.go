package operator

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
)

// Path wraps an AST node with its parent context, mirroring Babel's NodePath.
//
// The engine builds Path values during traversal; ParentPath is the chain up to
// the file root (nil for the root). Field/Index record where Node sits inside
// its Parent (a slice field such as BlockStmt.List, GenDecl.Specs or File.Decls
// uses Index >= 0; a scalar field keeps Index == -1).
type Path struct {
	Node       ast.Node
	Parent     ast.Node // the node that contains Node
	ParentPath *Path    // nil at file root
	Field      string   // e.g. "List", "Specs", "Decls"
	Index      int      // -1 when not in a slice field
}

// inRange reports whether idx is a valid index into a slice of length n.
func inRange(idx, n int) bool {
	return idx >= 0 && idx < n
}

// remove is the method form of Remove, mirroring putout's path.remove().
func (p *Path) remove() {
	Remove(p)
}

// Remove removes path.Node from its parent, preserving comments.
//
// It handles three parent shapes:
//   - *ast.BlockStmt  → splice List[path.Index]
//   - *ast.GenDecl    → splice Specs[path.Index] (for import removal)
//   - *ast.File       → splice Decls[path.Index] (for func/var removal)
//
// If path.Node is the last node in its parent slice and carries comment groups
// (e.g. a Doc comment), those comments are attached to the file's Comments
// slice so the printer keeps them after the node is gone. go/parser already
// records comments in File.Comments, so this only appends comment groups that
// were not already tracked, keeping the operation idempotent.
func Remove(path *Path) {
	switch p := path.Parent.(type) {
	case *ast.BlockStmt:
		if inRange(path.Index, len(p.List)) {
			if path.Index == len(p.List)-1 {
				preserveComments(path, p.List[path.Index])
			}
			p.List = append(p.List[:path.Index], p.List[path.Index+1:]...)
		}
	case *ast.GenDecl:
		if inRange(path.Index, len(p.Specs)) {
			if path.Index == len(p.Specs)-1 {
				preserveComments(path, p.Specs[path.Index])
			}
			p.Specs = append(p.Specs[:path.Index], p.Specs[path.Index+1:]...)
		}
	case *ast.File:
		if inRange(path.Index, len(p.Decls)) {
			if path.Index == len(p.Decls)-1 {
				preserveComments(path, p.Decls[path.Index])
			}
			p.Decls = append(p.Decls[:path.Index], p.Decls[path.Index+1:]...)
		}
		if inRange(path.Index, len(p.Imports)) {
			p.Imports = append(p.Imports[:path.Index], p.Imports[path.Index+1:]...)
		}
	}
}

// preserveComments appends any comment groups found within node to the file's
// Comments slice, skipping groups the file already tracks. This keeps comments
// attached to a removed last-in-scope node in the printer output without
// duplicating them.
func preserveComments(path *Path, node ast.Node) {
	groups := commentGroups(node)
	if len(groups) == 0 {
		return
	}
	file := findFile(path)
	if file == nil {
		return
	}
	for _, g := range groups {
		if slices.Contains(file.Comments, g) {
			continue
		}
		file.Comments = append(file.Comments, g)
	}
}

// findFile walks the ParentPath chain from path looking for the *ast.File root.
func findFile(path *Path) *ast.File {
	for cur := path; cur != nil; cur = cur.ParentPath {
		if f, ok := cur.Node.(*ast.File); ok {
			return f
		}
	}
	return nil
}

// commentGroups collects every non-nil *ast.CommentGroup reachable from node,
// descending through pointer and interface fields (Doc, Comment, etc).
func commentGroups(n ast.Node) []*ast.CommentGroup {
	if n == nil {
		return nil
	}
	v := reflect.ValueOf(n)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	t := v.Type()
	var groups []*ast.CommentGroup
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if t.Field(i).Type == reflect.TypeOf(&ast.CommentGroup{}) {
			if cg, ok := f.Interface().(*ast.CommentGroup); ok && cg != nil {
				groups = append(groups, cg)
			}
			continue
		}
		k := f.Kind()
		if k == reflect.Pointer || k == reflect.Interface {
			if f.IsNil() {
				continue
			}
			if child, ok := f.Interface().(ast.Node); ok {
				groups = append(groups, commentGroups(child)...)
			}
		}
	}
	return groups
}

// setPos sets every token.Pos field in node's sub-tree to pos, so the printer
// keeps the original source location while printing the replacement.
func setPos(n ast.Node, pos token.Pos) {
	ast.Inspect(n, func(node ast.Node) bool {
		if node == nil {
			return true
		}
		e := reflect.ValueOf(node).Elem()
		t := e.Type()
		for i := 0; i < e.NumField(); i++ {
			if t.Field(i).Type == reflect.TypeOf(token.NoPos) {
				e.Field(i).SetInt(int64(pos))
			}
		}
		return true
	})
}

// ReplaceWith replaces path.Node with node in path's parent, preserving the
// original node's position so the printer keeps its source location.
func ReplaceWith(path *Path, node ast.Node) {
	replaceInSlice(path, func(list reflect.Value, i int) {
		list.Set(reflect.AppendSlice(reflect.Append(list.Slice(0, i), reflect.ValueOf(node)), list.Slice(i+1, list.Len())))
	})
	setPos(node, path.Node.Pos())
}

// ReplaceWithMultiple replaces path.Node with several nodes in path's parent
// slice. Only valid when the parent field is a slice (List, Specs, Decls).
// The original node's position is preserved on nodes[0].
func ReplaceWithMultiple(path *Path, nodes []ast.Node) {
	if len(nodes) == 0 {
		return
	}
	setPos(nodes[0], path.Node.Pos())
	replaceSlice(path, nodes)
}

// InsertBefore inserts node before path.Node in path's parent slice.
func InsertBefore(path *Path, node ast.Node) {
	insertSlice(path, node, 0)
}

// InsertAfter inserts node after path.Node in path's parent slice.
func InsertAfter(path *Path, node ast.Node) {
	insertSlice(path, node, 1)
}

// parentSlice returns the slice field named by path.Field on the parent, or the
// zero Value when the parent is not a struct exposing that field.
func parentSlice(path *Path) reflect.Value {
	parent := reflect.ValueOf(path.Parent)
	if parent.Kind() != reflect.Pointer || parent.IsNil() || parent.Elem().Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return parent.Elem().FieldByName(path.Field)
}

// replaceInSlice applies apply to path's parent slice at path.Index. When the
// parent has no matching slice field or the index is out of range, it is a
// no-op.
func replaceInSlice(path *Path, apply func(list reflect.Value, i int)) {
	s := parentSlice(path)
	if !s.IsValid() || s.Kind() != reflect.Slice || !inRange(path.Index, s.Len()) {
		return
	}
	apply(s, path.Index)
}

// replaceSlice replaces the path.Index-th element of the parent slice with
// nodes.
func replaceSlice(path *Path, nodes []ast.Node) {
	replaceInSlice(path, func(s reflect.Value, i int) {
		var vals []reflect.Value
		for _, n := range nodes {
			vals = append(vals, reflect.ValueOf(n))
		}
		s.Set(reflect.AppendSlice(reflect.Append(s.Slice(0, i), vals...), s.Slice(i+1, s.Len())))
	})
}

// insertSlice inserts node in the parent slice at path.Index+offset.
func insertSlice(path *Path, node ast.Node, offset int) {
	replaceInSlice(path, func(s reflect.Value, i int) {
		s.Set(reflect.AppendSlice(reflect.Append(s.Slice(0, i+offset), reflect.ValueOf(node)), s.Slice(i+offset, s.Len())))
	})
}

// GetBinding returns the Path of the declaration of name visible from path,
// walking up the ParentPath chain. It returns nil when name is not declared.
// It covers:
//   - short variable declarations (:=)
//   - var declarations inside a block (DeclStmt)
//   - import specs
//   - top-level var / const / type declarations
//   - function declarations
//   - function parameters (FuncDecl and FuncLit)
func GetBinding(path *Path, name string) *Path {
	for cur := path; cur != nil; cur = cur.ParentPath {
		switch n := cur.Node.(type) {
		case *ast.BlockStmt:
			if stmt, ok := blockDeclaresStmt(n, name); ok {
				return &Path{Node: stmt, Parent: n}
			}
		case *ast.File:
			if node := fileDeclares(n, name); node != nil {
				return &Path{Node: node, Parent: n}
			}
		case *ast.FuncDecl:
			if paramsDeclare(n.Type.Params, name) {
				return &Path{Node: n, Parent: cur.Parent}
			}
		case *ast.FuncLit:
			if paramsDeclare(n.Type.Params, name) {
				return &Path{Node: n, Parent: cur.Parent}
			}
		}
	}
	return nil
}

// GetBindingPath is an alias for GetBinding (mirrors putout naming).
var GetBindingPath = GetBinding

// blockDeclaresStmt reports whether any statement in block declares name via a
// short variable declaration (:=) or a var declaration, returning the declaring
// statement.
func blockDeclaresStmt(block *ast.BlockStmt, name string) (ast.Stmt, bool) {
	for _, stmt := range block.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			if s.Tok != token.DEFINE {
				continue
			}
			for _, lhs := range s.Lhs {
				if id, ok := lhs.(*ast.Ident); ok && id.Name == name {
					return stmt, true
				}
			}
		case *ast.DeclStmt:
			gd, ok := s.Decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, n := range vs.Names {
					if n.Name == name {
						return stmt, true
					}
				}
			}
		}
	}
	return nil, false
}

// fileDeclares reports whether file declares name as an import, a top-level
// var/const/type, or a function, returning the declaring node.
func fileDeclares(file *ast.File, name string) ast.Node {
	for _, imp := range file.Imports {
		if imp.Name != nil && imp.Name.Name == name {
			return imp
		}
	}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name.Name == name {
				return decl
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					for _, n := range s.Names {
						if n.Name == name {
							return decl
						}
					}
				case *ast.TypeSpec:
					if s.Name.Name == name {
						return decl
					}
				}
			}
		}
	}
	return nil
}

// paramsDeclare reports whether fields declares name as a named parameter.
func paramsDeclare(fields *ast.FieldList, name string) bool {
	if fields == nil {
		return false
	}
	for _, f := range fields.List {
		for _, n := range f.Names {
			if n.Name == name {
				return true
			}
		}
	}
	return false
}

// Rename renames every *ast.Ident named from to to within the subtree rooted
// at path.Node (it does not cross out of that subtree).
func Rename(path *Path, from, to string) {
	ast.Inspect(path.Node, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok && id.Name == from {
			id.Name = to
		}
		return true
	})
}
