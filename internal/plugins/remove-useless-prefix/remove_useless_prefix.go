package remove_useless_prefix

import (
	"go/ast"
	"reflect"

	. "coderaiser/indra/types"
)

const goTapePath = `"github.com/coderaiser/go-tape"`

func Report() string { return "remove useless tape prefix" }

func Traverse() Traverser {
	return Traverser{"*ast.File": visitFile}
}

func visitFile(node ast.Node, _ Vars) []Place {
	file := node.(*ast.File)
	alias, _ := findTapeImport(file)
	if alias == "" {
		return nil
	}
	return []Place{{Message: Report()}}
}

// findTapeImport returns the local alias (and its import spec) for a named
// (non-dot, non-blank) import of the go-tape package. A dot or blank import
// returns an empty alias.
func findTapeImport(file *ast.File) (string, *ast.ImportSpec) {
	for _, imp := range file.Imports {
		if imp.Path.Value != goTapePath {
			continue
		}
		if imp.Name == nil || imp.Name.Name == "." || imp.Name.Name == "_" {
			return "", nil
		}
		return imp.Name.Name, imp
	}
	return "", nil
}

// Fix rewrites a named go-tape import to a dot import and drops the alias
// prefix from every selector use (tape.X → X).
func Fix(node ast.Node, _ []Place) {
	file := node.(*ast.File)
	alias, spec := findTapeImport(file)
	if alias == "" {
		return
	}
	spec.Name = &ast.Ident{Name: ".", NamePos: spec.Name.NamePos}
	replaceSelectors(reflect.ValueOf(file), alias)
}

// typed variants used to prune the semantic (non-syntactic) parts of the tree,
// which can contain back-references (Ident.Obj -> Decl -> Ident) and produce
// reflection cycles if walked.
var (
	astIdentType = reflect.TypeOf((*ast.Ident)(nil))
	astScopeType = reflect.TypeOf((*ast.Scope)(nil))
)

// replaceSelectors walks an AST value via reflection, replacing every
// *ast.SelectorExpr whose X is an *ast.Ident named alias with just its Sel
// ident in the containing settable field.
func replaceSelectors(v reflect.Value, alias string) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return
		}
		elem := v.Elem()
		if sel, ok := elem.Interface().(*ast.SelectorExpr); ok && selMatchesAlias(sel, alias) {
			if v.CanSet() {
				v.Set(reflect.ValueOf(sel.Sel))
			}
			return
		}
		replaceSelectors(elem, alias)
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		// Skip identifiers and scopes: identifiers hold Obj back-references and
		// scopes are semantic, neither can contain a selector to rewrite.
		if v.Type() == astIdentType || v.Type() == astScopeType {
			return
		}
		replaceInStruct(v.Elem(), alias)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if maybeReplaceSel(elem, alias) {
				continue
			}
			replaceSelectors(elem, alias)
		}
	}
}

// replaceInStruct walks every settable field of a struct value.
func replaceInStruct(v reflect.Value, alias string) {
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		if maybeReplaceSel(f, alias) {
			continue
		}
		switch f.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Slice:
			replaceSelectors(f, alias)
		}
	}
}

// selMatchesAlias reports whether sel is a selector whose X is an *ast.Ident
// named alias.
func selMatchesAlias(sel *ast.SelectorExpr, alias string) bool {
	ident, ok := sel.X.(*ast.Ident)
	return ok && ident.Name == alias
}

// maybeReplaceSel checks if a pointer- or interface-typed value holds an
// *ast.SelectorExpr whose X matches the alias. If the value is settable it is
// replaced with the Sel ident.
func maybeReplaceSel(f reflect.Value, alias string) bool {
	var iface interface{}
	switch f.Kind() {
	case reflect.Interface:
		if f.IsNil() {
			return false
		}
		iface = f.Interface()
	case reflect.Ptr:
		if f.IsNil() {
			return false
		}
		iface = f.Interface()
	default:
		return false
	}
	sel, ok := iface.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if !selMatchesAlias(sel, alias) {
		return false
	}
	if f.CanSet() {
		f.Set(reflect.ValueOf(sel.Sel))
	}
	return true
}
