// Package babel mirrors @putout/babel: type-check helpers plus a traverse
// bootstrapper and Path navigation, so external consumers (e.g. go-printer)
// can walk and inspect Go ASTs through a single import.
//
// Indra's own plugins keep dot-importing operator and types; babel carries
// its own Path type because Go cannot define methods on an alias of a
// foreign package's type.
package babel

import (
	"go/ast"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"
)

// Path is a node together with its ancestor stack and apply cursor. It mirrors
// types.Path field-for-field (convertible with babel.Path(types.Path{...})),
// but is a distinct type so the Get / GetAll / Type / IsEmpty methods below
// can live in this package.
type Path struct {
	Node   ast.Node
	Stack  []ast.Node // ancestors root-first, excluding Node
	Cursor *astutil.Cursor
}

// Traverse bootstraps a root Path from node and walks it with visitors,
// routing each node to the matching visitor by type name (e.g. "CallExpr").
// Mirrors @putout/babel traverse — used by go-printer to start a walk from a
// raw ast.Node without a pre-existing Path.
func Traverse(root ast.Node, visitors map[string]func(Path)) {
	astutil.Apply(root, func(c *astutil.Cursor) bool {
		node := c.Node()
		if node == nil {
			return true
		}
		typeName := reflect.TypeOf(node).Elem().Name()
		if fn, ok := visitors[typeName]; ok {
			fn(Path{Node: node, Cursor: c})
		}
		return true
	}, nil)
}

// Get returns the child Path at the named struct field of p.Node.
// Field is the Go struct field name: "Fun", "Args", "Body", "X", etc.
// Returns a zero Path (IsEmpty() == true) if p.Node is nil, the field does
// not exist, or is not a non-nil single ast.Node.
// Mirrors Babel's path.get(field) for single-node fields.
func (p Path) Get(field string) Path {
	if p.Node == nil {
		return Path{}
	}
	v := reflect.ValueOf(p.Node)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return Path{}
		}
		v = v.Elem()
	}
	f := v.FieldByName(field)
	if !f.IsValid() {
		return Path{}
	}
	node, ok := f.Interface().(ast.Node)
	if !ok || reflect.ValueOf(node).IsNil() {
		return Path{}
	}
	return Path{Node: node}
}

// GetAll returns child Paths for an array field of p.Node.
// Field is the Go struct field name: "Args", "Elts", "List", "Specs", etc.
// Returns nil if p.Node is nil, the field does not exist, or is not a slice;
// elements that are not ast.Node are skipped.
// Mirrors Babel's path.get(field) for array fields.
func (p Path) GetAll(field string) []Path {
	if p.Node == nil {
		return nil
	}
	v := reflect.ValueOf(p.Node)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Slice {
		return nil
	}
	paths := make([]Path, 0, f.Len())
	for i := range f.Len() {
		elem := f.Index(i)
		node, ok := elem.Interface().(ast.Node)
		if !ok {
			continue
		}
		paths = append(paths, Path{Node: node})
	}
	return paths
}

// Type returns the Go type name of p.Node without package prefix,
// e.g. "CallExpr", "FuncDecl", "Ident". Returns "" for a zero Path.
// Mirrors Babel's path.type.
func (p Path) Type() string {
	if p.Node == nil {
		return ""
	}
	t := reflect.TypeOf(p.Node)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

// IsEmpty reports whether p.Node is nil.
// Mirrors Babel's !path.node check in printer.
func (p Path) IsEmpty() bool {
	return p.Node == nil
}
