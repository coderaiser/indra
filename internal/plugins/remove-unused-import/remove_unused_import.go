package removeunusedimport

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"coderaiser/indra/internal/engine"
)

var Plugin = engine.Plugin{
	Name:   "remove-unused-import",
	Report: func() string { return "remove unused import" },
	Traverse: func() map[string]engine.TraverseVisitor {
		return map[string]engine.TraverseVisitor{
			"*ast.File": visitFile,
		}
	},
}

func visitFile(node ast.Node, _ engine.Vars) []engine.Place {
	file := node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(file)

	var places []engine.Place
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			places = append(places, engine.Place{
				Message: fmt.Sprintf("remove unused import: %s", imp.path),
			})
		}
	}
	return places
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
				info.localName = filepath.Base(raw)
			}
			imports = append(imports, info)
		}
	}
	return imports
}

func countIdentUses(file *ast.File) map[string]int {
	used := make(map[string]int)
	ast.Inspect(file, func(n ast.Node) bool {
		// skip import decls themselves
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			return false
		}
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		used[ident.Name]++
		return true
	})
	return used
}
