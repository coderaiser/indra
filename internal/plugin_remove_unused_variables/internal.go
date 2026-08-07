package remove_unused_variables

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// ── imports ──────────────────────────────────────────────────────────────────

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

// importFinding is a synthetic ast.Node that pairs a file with one unused
// import spec so Report and Fix can act on the specific spec.
type importFinding struct {
	file *ast.File
	spec *ast.ImportSpec
}

func (f *importFinding) Pos() token.Pos { return f.spec.Pos() }
func (f *importFinding) End() token.Pos { return f.spec.End() }

// findUnusedImportsAndConsts visits *ast.File nodes during traversal.
// It calls push once per unused import and once for unused consts.
func findUnusedImportsAndConsts(node ast.Node, push func(ast.Node)) {
	file := node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(file)
	found := false
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] == 0 {
			push(&importFinding{file: file, spec: imp.spec})
			found = true
		}
	}
	if !found && len(unusedConstNames(file)) > 0 {
		push(file)
	}
}

// fixOneUnusedImport removes a single import spec from the file's AST.
func fixOneUnusedImport(file *ast.File, target *ast.ImportSpec) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		kept := genDecl.Specs[:0]
		for _, spec := range genDecl.Specs {
			if spec != target {
				kept = append(kept, spec)
			}
		}
		genDecl.Specs = kept
	}
	// drop empty import blocks
	kept := file.Decls[:0]
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT && len(genDecl.Specs) == 0 {
			continue
		}
		kept = append(kept, decl)
	}
	file.Decls = kept
	// sync file.Imports
	newImports := file.Imports[:0]
	for _, imp := range file.Imports {
		if imp != target {
			newImports = append(newImports, imp)
		}
	}
	file.Imports = newImports
}

// ── consts ───────────────────────────────────────────────────────────────────

// unusedConstNames returns names of file-level consts never referenced.
func unusedConstNames(file *ast.File) []string {
	declared := collectConstNames(file)
	if len(declared) == 0 {
		return nil
	}
	used := countIdentUses(file) // reuse existing; counts all idents outside import blocks
	var unused []string
	for _, name := range declared {
		if used[name] <= 1 { // 1 = the const name appears at its own definition
			unused = append(unused, name)
		}
	}
	return unused
}

func collectConstNames(file *ast.File) []string {
	var names []string
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range vs.Names {
				if name.Name != "_" {
					names = append(names, name.Name)
				}
			}
		}
	}
	return names
}

func fixUnusedConsts(file *ast.File) {
	unused := make(map[string]bool)
	for _, name := range unusedConstNames(file) {
		unused[name] = true
	}
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		kept := genDecl.Specs[:0]
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				kept = append(kept, spec)
				continue
			}
			allUnused := true
			for _, name := range vs.Names {
				if !unused[name.Name] {
					allUnused = false
				}
			}
			if !allUnused {
				kept = append(kept, spec)
			}
		}
		genDecl.Specs = kept
	}
	// drop empty const blocks
	kept := file.Decls[:0]
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.CONST && len(genDecl.Specs) == 0 {
			continue
		}
		kept = append(kept, decl)
	}
	file.Decls = kept
}

// ── variables ────────────────────────────────────────────────────────────────

// unusedVarNames returns the names declared via `:=` or `var` in block that
// are never read afterwards.
func unusedVarNames(block *ast.BlockStmt) []string {
	var decls []string
	seen := map[string]bool{}

	for _, stmt := range block.List {
		// handle: x := expr
		if assign, ok := stmt.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
			for _, lhs := range assign.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident.Name == "_" {
					continue
				}
				if !seen[ident.Name] {
					seen[ident.Name] = true
					decls = append(decls, ident.Name)
				}
			}
			continue
		}
		// handle: var x = expr  or  var x int
		if declStmt, ok := stmt.(*ast.DeclStmt); ok {
			genDecl, ok := declStmt.Decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, name := range vs.Names {
					if name.Name == "_" {
						continue
					}
					if !seen[name.Name] {
						seen[name.Name] = true
						decls = append(decls, name.Name)
					}
				}
			}
		}
	}

	if len(decls) == 0 {
		return nil
	}

	reads := map[string]int{}
	for _, stmt := range block.List {
		// :=  — only count RHS, not the declared names
		if assign, ok := stmt.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
			for _, rhs := range assign.Rhs {
				countIdents(rhs, reads)
			}
			continue
		}
		// var x = expr — only count RHS values, not the declared names
		if declStmt, ok := stmt.(*ast.DeclStmt); ok {
			if genDecl, ok := declStmt.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for _, val := range vs.Values {
							countIdents(val, reads)
						}
					}
				}
				continue
			}
		}
		countIdents(stmt, reads)
	}

	var unused []string
	for _, d := range decls {
		if reads[d] == 0 {
			unused = append(unused, d)
		}
	}
	return unused
}

func countIdents(n ast.Node, reads map[string]int) {
	ast.Inspect(n, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		if ok {
			reads[ident.Name]++
		}
		return true
	})
}

// findUnusedVars visits *ast.BlockStmt nodes during traversal.
func findUnusedVars(node ast.Node, push func(ast.Node)) {
	block := node.(*ast.BlockStmt)
	if len(unusedVarNames(block)) > 0 {
		push(block)
	}
}

// fixUnusedVars removes unused variables from a block in place.
func fixUnusedVars(block *ast.BlockStmt) {
	unused := make(map[string]bool, len(unusedVarNames(block)))
	for _, n := range unusedVarNames(block) {
		unused[n] = true
	}
	kept := block.List[:0]
	for _, stmt := range block.List {
		// handle: x := expr
		if assign, ok := stmt.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
			allUnused := true
			for _, lhs := range assign.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident.Name == "_" {
					continue
				}
				if !unused[ident.Name] {
					allUnused = false
				}
			}
			if allUnused {
				continue // drop statement
			}
			for _, lhs := range assign.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if ok && unused[ident.Name] {
					ident.Name = "_"
				}
			}
			kept = append(kept, stmt)
			continue
		}
		// handle: var x = expr
		if declStmt, ok := stmt.(*ast.DeclStmt); ok {
			genDecl, ok := declStmt.Decl.(*ast.GenDecl)
			if ok && genDecl.Tok == token.VAR {
				newSpecs := genDecl.Specs[:0]
				for _, spec := range genDecl.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						newSpecs = append(newSpecs, spec)
						continue
					}
					allUnused := true
					for _, name := range vs.Names {
						if name.Name != "_" && !unused[name.Name] {
							allUnused = false
						}
					}
					if !allUnused {
						newSpecs = append(newSpecs, spec)
					}
				}
				genDecl.Specs = newSpecs
				if len(genDecl.Specs) > 0 {
					kept = append(kept, stmt)
				}
				continue
			}
		}
		kept = append(kept, stmt)
	}
	block.List = kept
}
