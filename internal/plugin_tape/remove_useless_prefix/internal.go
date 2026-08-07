package remove_useless_prefix

import (
	"go/ast"
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

// selMatchesAlias reports whether sel is a selector whose X is an *ast.Ident
// named alias.
func selMatchesAlias(sel *ast.SelectorExpr, alias string) bool {
	ident, ok := sel.X.(*ast.Ident)
	return ok && ident.Name == alias
}
