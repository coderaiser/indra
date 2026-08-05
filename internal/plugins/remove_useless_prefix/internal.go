package remove_useless_prefix

import (
	"go/ast"
	"reflect"
)

const goTapePath = `"github.com/coderaiser/go-tape"`

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

// usedBareNames returns all bare *ast.Ident names that appear outside any
// selector expression — i.e. neither as a qualifier (X) nor as a member (Sel)
// of a SelectorExpr. These are names that would clash if a prefix were removed
// and a same-named selector member became bare: removing the prefix turns an
// alias.X selector into a bare X, which collides with any pre-existing bare X.
func usedBareNames(file *ast.File) map[string]bool {
	names := make(map[string]bool)
	ast.Inspect(file, func(n ast.Node) bool {
		if _, ok := n.(*ast.SelectorExpr); ok {
			return false // skip the whole selector (X qualifier and Sel member)
		}
		if id, ok := n.(*ast.Ident); ok {
			names[id.Name] = true
		}
		return true
	})
	return names
}

// hasLocalCollision reports whether removing the alias prefix from any
// alias.X selector would introduce an identifier that collides with a
// locally declared (package-level) name or with a name already used as a bare
// identifier elsewhere. Skipping such files avoids emitting broken code where
// the unqualified name would now resolve to a local decl or clash with an
// existing bare usage.
func hasLocalCollision(file *ast.File, alias string) bool {
	declared := declaredNames(file)
	usedBare := usedBareNames(file)
	collision := false
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if !selMatchesAlias(sel, alias) {
			return true
		}
		if declared[sel.Sel.Name] || usedBare[sel.Sel.Name] {
			collision = true
			return false
		}
		return true
	})
	return collision
}

// declaredNames returns the set of package-level declared identifiers.
func declaredNames(file *ast.File) map[string]bool {
	names := make(map[string]bool)
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name != nil {
				names[d.Name.Name] = true
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name != nil {
						names[s.Name.Name] = true
					}
				case *ast.ValueSpec:
					for _, n := range s.Names {
						names[n.Name] = true
					}
				}
			}
		}
	}
	return names
}

// astIdentType is used to prune *ast.Ident during reflection walks: identifiers
// hold Obj back-references that can form cycles if walked.
// astScopeName identifies ast.Scope by reflection without referencing the
// deprecated type directly — scopes are semantic and must also be skipped to
// avoid cycles.
var (
	astIdentType    = reflect.TypeOf((*ast.Ident)(nil))
	astScopePkgPath = "go/ast"
	astScopeName    = "Scope"
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
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		// Skip identifiers and scopes: identifiers hold Obj back-references and
		// scopes are semantic, neither can contain a selector to rewrite.
		if v.Type() == astIdentType {
			return
		}
		// Identify ast.Scope by package path and name to avoid referencing the
		// deprecated type directly.
		if elem := v.Type().Elem(); elem.PkgPath() == astScopePkgPath && elem.Name() == astScopeName {
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
		case reflect.Interface, reflect.Pointer, reflect.Slice:
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
	case reflect.Pointer:
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
