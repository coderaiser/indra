package remove_unused_import

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

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

func findUnusedImports(node ast.Node, push func(ast.Node)) {
	file := node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(file)

	hasUnused := false
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			hasUnused = true
			break
		}
	}
	if hasUnused {
		push(file)
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

type importInfo struct {
	spec      *ast.ImportSpec
	localName string
	path      string
	blank     bool
	dot       bool
}

func collectImports(file *ast.File) []importInfo {
	var imports []importInfo
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			info := importInfo{
				spec: importSpec,
				path: importSpec.Path.Value,
			}
			if importSpec.Name != nil {
				switch importSpec.Name.Name {
				case "_":
					info.blank = true
				case ".":
					info.dot = true
				default:
					info.localName = importSpec.Name.Name
				}
			} else {
				raw := strings.Trim(importSpec.Path.Value, `"`)
				// Go derives the implicit local name from the package's declared
				// name. Path basenames containing hyphens (e.g. add-t-end) are
				// declared as underscore names (add_t_end), so match that form.
				info.localName = strings.ReplaceAll(filepath.Base(raw), "-", "_")
			}
			imports = append(imports, info)
		}
	}
	return imports
}

func countIdentUses(file *ast.File) map[string]int {
	used := make(map[string]int)
	ast.Inspect(file, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			return false
		}
		ident, ok := n.(*ast.Ident)
		// A hand-built *ast.File may have a typed-nil Name (ast.Walk visits
		// File.Name unconditionally), so guard against a nil Ident.
		if !ok || ident == nil {
			return true
		}
		used[ident.Name]++
		return true
	})
	return used
}
