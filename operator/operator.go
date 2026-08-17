package operator

import (
	"go/ast"
	"go/token"
	"reflect"
	"slices"
	"strconv"

	"coderaiser/indra/types"
)

// Extract returns the string value of node:
//   - *ast.Ident        → Name
//   - *ast.BasicLit     → Value (unquoted for strings)
//   - *ast.SelectorExpr → "X.Sel"
//
// It panics on unsupported types (same as putout's extract).
func Extract(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.BasicLit:
		if n.Kind == token.STRING {
			if s, err := strconv.Unquote(n.Value); err == nil {
				return s
			}
		}
		return n.Value
	case *ast.SelectorExpr:
		return Extract(n.X) + "." + n.Sel.Name
	default:
		panic("operator: Extract: unsupported node type")
	}
}

// IsSimple reports whether node is a literal, identifier, or selector
// expression — i.e. has no sub-expressions that could be extracted.
func IsSimple(node ast.Node) bool {
	switch node.(type) {
	case *ast.BasicLit, *ast.Ident, *ast.SelectorExpr:
		return true
	}
	return false
}


// Remove deletes path.Node from its parent, preserving comments by migrating
// them to the file's Comments slice when path.Node is the last child in its
// parent.
//
// Mirrors putout's operator remove(path).
func Remove(path types.Path) {
	preserveComments(path)
	path.Delete()
}

// ReplaceWith replaces path.Node with node in path's parent, preserving the
// original node's position so the printer keeps its source location.
//
// Mirrors putout's operator replaceWith(path, node).
func ReplaceWith(path types.Path, node ast.Node) {
	copyPos(node, path.Node.Pos())
	path.Replace(node)
}

// ReplaceWithMultiple replaces path.Node with multiple nodes in path's parent
// slice. Preserves the original node's position on nodes[0].
func ReplaceWithMultiple(path types.Path, nodes []ast.Node) {
	if len(nodes) == 0 {
		return
	}
	copyPos(nodes[0], path.Node.Pos())
	for index := len(nodes) - 1; index >= 1; index-- {
		path.InsertAfter(nodes[index])
	}
	path.Replace(nodes[0])
}

// InsertBefore inserts node before path.Node in path's parent slice.
// Mirrors putout's operator insertBefore(path, node).
func InsertBefore(path types.Path, node ast.Node) {
	path.InsertBefore(node)
}

// InsertAfter inserts node after path.Node in path's parent slice,
// preserving trailing comments of path.Node.
// Mirrors putout's operator insertAfter(path, node).
func InsertAfter(path types.Path, node ast.Node) {
	path.InsertAfter(node)
}

// GetBinding returns the Path of the declaration of name visible from path,
// walking up the Stack. Returns nil when name is not declared.
//
// Scanned declaration forms:
//   - *ast.BlockStmt  → short variable declarations (:=), var declarations
//   - *ast.FuncDecl   → receiver, parameters, named return values
//   - *ast.FuncLit    → parameters, named return values
//   - *ast.File       → import specs, top-level funcs/vars/types
//
// Mirrors putout's operator getBinding(path, name).
func GetBinding(path types.Path, name string) *types.Path {
	for index := len(path.Stack) - 1; index >= 0; index-- {
		switch ancestor := path.Stack[index].(type) {
		case *ast.BlockStmt:
			if stmt, ok := blockDeclaresStmt(ancestor, name); ok {
				result := types.Path{Node: stmt}
				return &result
			}
		case *ast.FuncDecl:
			if paramsDeclare(ancestor.Type.Params, name) {
				result := types.Path{Node: ancestor}
				return &result
			}
		case *ast.FuncLit:
			if paramsDeclare(ancestor.Type.Params, name) {
				result := types.Path{Node: ancestor}
				return &result
			}
		case *ast.File:
			if node := fileDeclares(ancestor, name); node != nil {
				result := types.Path{Node: node}
				return &result
			}
		}
	}
	return nil
}

// GetBindingPath is an alias for GetBinding. Mirrors putout naming.
var GetBindingPath = GetBinding


// preserveComments appends any comment groups found within path.Node to the
// file's Comments slice, skipping groups the file already tracks.
func preserveComments(path types.Path) {
	groups := commentGroups(path.Node)
	if len(groups) == 0 {
		return
	}
	file := findFile(path)
	if file == nil {
		return
	}
	for _, group := range groups {
		if slices.Contains(file.Comments, group) {
			continue
		}
		file.Comments = append(file.Comments, group)
	}
}

// findFile walks path.Stack looking for the *ast.File root.
func findFile(path types.Path) *ast.File {
	for index := len(path.Stack) - 1; index >= 0; index-- {
		if file, ok := path.Stack[index].(*ast.File); ok {
			return file
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

// copyPos sets every token.Pos field in node's sub-tree to pos, so the printer
// keeps the original source location while printing the replacement.
func copyPos(n ast.Node, pos token.Pos) {
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

// Rename renames all *ast.Ident nodes named from to to within the subtree
// rooted at path.Node. Does not cross function literal boundaries.
// Mirrors putout's operator rename(path, from, to).
func Rename(path types.Path, from, to string) {
	ast.Inspect(path.Node, func(node ast.Node) bool {
		if identifier, ok := node.(*ast.Ident); ok && identifier.Name == from {
			identifier.Name = to
		}
		return true
	})
}

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

