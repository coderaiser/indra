package rules

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"coderaiser/indra/internal/lint/rule"
)

type RemoveUnusedImport struct{}

func (r *RemoveUnusedImport) Name() string {
	return "remove-unused-import"
}

func (r *RemoveUnusedImport) Check(file *ast.File, fset *token.FileSet) []rule.Result {
	var results []rule.Result

	imports := collectImports(file)
	used := countIdentUses(file)

	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			results = append(results, rule.Result{
				Pos:     fset.Position(imp.spec.Pos()),
				Message: fmt.Sprintf("remove unused import: %s", imp.path),
			})
		}
	}

	return results
}

func (r *RemoveUnusedImport) Fix(file *ast.File, fset *token.FileSet) bool {
	imports := collectImports(file)
	used := countIdentUses(file)

	modified := false

	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			removeImportSpec(file, imp)
			modified = true
		}
	}

	return modified
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
				if importSpec.Name.Name == "_" {
					info.blank = true
				} else if importSpec.Name.Name == "." {
					info.dot = true
				} else {
					info.localName = importSpec.Name.Name
				}
			} else {
				raw := strings.Trim(importSpec.Path.Value, "\"")
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
		genDecl, ok := n.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			return false
		}

		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name == "_" || ident.Name == "." {
			return true
		}
		used[ident.Name]++
		return true
	})

	return used
}

func removeImportSpec(file *ast.File, imp importInfo) {
	for i, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}

		for j, spec := range genDecl.Specs {
			if spec == imp.spec {
				genDecl.Specs = append(genDecl.Specs[:j], genDecl.Specs[j+1:]...)
				if len(genDecl.Specs) == 0 {
					file.Decls = append(file.Decls[:i], file.Decls[i+1:]...)
				}
				return
			}
		}
	}
}