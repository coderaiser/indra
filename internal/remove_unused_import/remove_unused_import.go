package remove_unused_import

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/types"
)

func Report(node ast.Node) string {
	if node == nil {
		return "remove unused import"
	}
	file := node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(file)
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			return "remove unused import: " + imp.path
		}
	}
	return "remove unused import"
}

func Traverse() Traverser {
	return Traverser{
		"*ast.File": findUnusedImports,
	}
}

// Fix removes unused imports from the AST in place.
// node is *ast.File; options is unused.
func Fix(node ast.Node, _ map[string]any) {
	file := node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(file)

	unused := make(map[string]bool)
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			unused[imp.path] = true
		}
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		kept := genDecl.Specs[:0]
		for _, spec := range genDecl.Specs {
			imp := spec.(*ast.ImportSpec)
			if !unused[imp.Path.Value] {
				kept = append(kept, spec)
			}
		}
		genDecl.Specs = kept
	}
	kept := file.Decls[:0]
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT && len(genDecl.Specs) == 0 {
			continue
		}
		kept = append(kept, decl)
	}
	file.Decls = kept
}
