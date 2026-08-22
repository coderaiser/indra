package remove_useless_prefix

import (
	"go/ast"

	. "coderaiser/indra/operator"
	. "coderaiser/indra/types"
)

const goTapePath = `"github.com/coderaiser/go-tape"`

func Report(_ Path) string { return "remove useless tape prefix" }

// Fix rewrites a named go-tape import to a dot import and drops the alias
// prefix from every selector use (tape.X → X). It is only ever invoked on a
// pushed file, so findUselessPrefix has already guaranteed a non-empty alias
// with no local collision.
func Fix(p Path, _ map[string]any) {
	file := p.Node.(*ast.File)
	alias, spec := findTapeImport(file)

	spec.Name = &ast.Ident{Name: ".", NamePos: spec.Name.NamePos}

	p.Traverse(map[string]func(Path){
		"*ast.SelectorExpr": func(sp Path) {
			sel, ok := sp.Node.(*ast.SelectorExpr)
			if !ok || !selMatchesAlias(sel, alias) {
				return
			}
			ReplaceWith(sp, sel.Sel)
		},
	})
}

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessPrefix}
}

func findUselessPrefix(p Path, push func(Path)) {
	file := p.Node.(*ast.File)
	alias, _ := findTapeImport(file)
	if alias == "" {
		return
	}
	if hasLocalCollision(p, alias) {
		return
	}
	push(p)
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

// usedBareNames returns all bare *ast.Ident names that appear outside any
// selector expression — i.e. neither as a qualifier (X) nor as a member (Sel)
// of a SelectorExpr. These are names that would clash if a prefix were removed
// and a same-named selector member became bare: removing the prefix turns an
// alias.X selector into a bare X, which collides with any pre-existing bare X.
func usedBareNames(p Path) map[string]bool {
	names := make(map[string]bool)
	p.Traverse(map[string]func(Path){
		"*ast.SelectorExpr": func(sp Path) { sp.Skip() }, // skip the whole selector
		"*ast.Ident": func(ip Path) {
			names[ip.Node.(*ast.Ident).Name] = true
		},
	})
	return names
}

// hasLocalCollision reports whether removing the alias prefix from any
// alias.X selector would introduce an identifier that collides with a
// locally declared (package-level) name or with a name already used as a bare
// identifier elsewhere. Skipping such files avoids emitting broken code where
// the unqualified name would now resolve to a local decl or clash with an
// existing bare usage.
func hasLocalCollision(p Path, alias string) bool {
	file := p.Node.(*ast.File)
	declared := declaredNames(file)
	usedBare := usedBareNames(p)
	collision := false
	p.Traverse(map[string]func(Path){
		"*ast.SelectorExpr": func(sp Path) {
			sel := sp.Node.(*ast.SelectorExpr)
			if !selMatchesAlias(sel, alias) {
				return
			}
			if declared[sel.Sel.Name] || usedBare[sel.Sel.Name] {
				collision = true
				sp.Stop()
			}
		},
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

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
