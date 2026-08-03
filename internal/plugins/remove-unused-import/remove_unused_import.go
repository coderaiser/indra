package remove_unused_import

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	. "coderaiser/indra/types"
)

func Report() string { return "remove unused import" }

func Traverse() Traverser {
	return Traverser{
		"*ast.File": visitFile,
	}
}

func visitFile(node ast.Node, _ Vars) []Place {
	file := node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(file)

	var places []Place
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			places = append(places, Place{
				Message: fmt.Sprintf("remove unused import: %s", imp.path),
			})
		}
	}
	return places
}

// Fix removes unused imports from the AST in place.
// node is *ast.File. places contains findings from Traverse.
func Fix(node ast.Node, places []Place) {
	file := node.(*ast.File)
	unused := make(map[string]bool, len(places))
	for _, p := range places {
		unused[strings.TrimPrefix(p.Message, "remove unused import: ")] = true
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
