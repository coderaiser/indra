// Package remove_unused_variables reports file-scoped declarations that are
// never referenced within the same file: imports, consts, variables, and
// unexported functions. The engine lints one file at a time, so a helper that
// this file uses only from a sibling file would look unused here; keeping the
// whole implementation in one file keeps that false positive out of reach.
package remove_unused_variables

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"coderaiser/indra/types"
)

func Report(p types.Path) string {
	switch n := p.Node.(type) {
	case *importFinding:
		return "remove unused import: " + n.spec.Path.Value
	case *funcDeclFinding:
		return "remove unused private function: " + n.decl.Name.Name
	case *ast.File:
		consts := unusedConstNames(n)
		if len(consts) > 0 {
			return "remove unused const: " + consts[0]
		}
	case *ast.BlockStmt:
		unused := unusedVarNames(n)
		if len(unused) > 0 {
			return "remove unused variable: " + unused[0]
		}
	}
	return "remove unused variable"
}

func Traverse() types.Traverser {
	return types.Traverser{
		"*ast.File":      findUnusedImportsAndConsts,
		"*ast.BlockStmt": findUnusedVars,
	}
}

func Fix(p types.Path, _ map[string]any) {
	switch n := p.Node.(type) {
	case *importFinding:
		fixOneUnusedImport(n.file, n.spec)
	case *funcDeclFinding:
		fixOneUnusedPrivateFunc(n.file, n.decl)
	case *ast.File:
		fixUnusedConsts(n)
	case *ast.BlockStmt:
		fixUnusedVars(n)
	}
}

type Plugin struct{}

func (Plugin) Report(p types.Path) string            { return Report(p) }
func (Plugin) Traverse() types.Traverser             { return Traverse() }
func (Plugin) Fix(p types.Path, opts map[string]any) { Fix(p, opts) }

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
			importSpec := spec.(*ast.ImportSpec)
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

func countIdentUses(p types.Path) map[string]int {
	used := make(map[string]int)
	p.Traverse(map[string]func(types.Path){
		"*ast.GenDecl": func(gp types.Path) {
			if gp.Node.(*ast.GenDecl).Tok == token.IMPORT {
				gp.Skip()
			}
		},
		"*ast.Ident": func(ip types.Path) {
			id := ip.Node.(*ast.Ident)
			if id != nil {
				used[id.Name]++
			}
		},
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

// funcDeclFinding pairs a file with one unused private function declaration so
// Report and Fix can act on the specific decl, mirroring importFinding.
type funcDeclFinding struct {
	file *ast.File
	decl *ast.FuncDecl
}

func (f *funcDeclFinding) Pos() token.Pos { return f.decl.Pos() }
func (f *funcDeclFinding) End() token.Pos { return f.decl.End() }

// collectSelectorQualifiers returns all X identifiers from X.Y selector
// expressions in the file, excluding the import block.
func collectSelectorQualifiers(filePath types.Path) map[string]bool {
	qualifiers := make(map[string]bool)
	filePath.Traverse(map[string]func(types.Path){
		"*ast.GenDecl": func(declPath types.Path) {
			genDecl := declPath.Node.(*ast.GenDecl)
			if genDecl.Tok == token.IMPORT {
				declPath.Skip()
			}
		},
		"*ast.SelectorExpr": func(selectorPath types.Path) {
			selector := selectorPath.Node.(*ast.SelectorExpr)
			if ident, ok := selector.X.(*ast.Ident); ok {
				qualifiers[ident.Name] = true
			}
		},
	})
	return qualifiers
}

// collectDeclaredNames returns all package-level declared identifiers.
func collectDeclaredNames(filePath types.Path) map[string]bool {
	names := make(map[string]bool)
	filePath.Traverse(map[string]func(types.Path){
		"*ast.FuncDecl": func(funcPath types.Path) {
			funcDecl := funcPath.Node.(*ast.FuncDecl)
			if funcDecl.Name != nil {
				names[funcDecl.Name.Name] = true
			}
			funcPath.Skip()
		},
		"*ast.TypeSpec": func(typePath types.Path) {
			typeSpec := typePath.Node.(*ast.TypeSpec)
			if typeSpec.Name != nil {
				names[typeSpec.Name.Name] = true
			}
		},
		"*ast.ValueSpec": func(valuePath types.Path) {
			valueSpec := valuePath.Node.(*ast.ValueSpec)
			for _, name := range valueSpec.Names {
				names[name.Name] = true
			}
		},
	})
	return names
}

// findUnusedImportsAndConsts visits *ast.File nodes during traversal.
// It calls push once per unused import and once for unused consts.
func findUnusedImportsAndConsts(p types.Path, push func(types.Path)) {
	file := p.Node.(*ast.File)
	imports := collectImports(file)
	used := countIdentUses(p)
	qualifiers := collectSelectorQualifiers(p)
	declared := collectDeclaredNames(p)

	// accounted: names explained by imports or local declarations
	accounted := make(map[string]bool)
	for _, imp := range imports {
		if !imp.blank && !imp.dot {
			accounted[imp.localName] = true
		}
	}
	for name := range declared {
		accounted[name] = true
	}

	// unaccounted qualifiers: used in selectors but not explained
	var unaccounted []string
	for qualifier := range qualifiers {
		if !accounted[qualifier] {
			unaccounted = append(unaccounted, qualifier)
		}
	}

	// ambiguous: unnamed imports whose basename-derived name is not in used
	var ambiguous []importInfo
	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] > 0 {
			continue
		}
		ambiguous = append(ambiguous, imp)
	}

	// pair ambiguous imports with unaccounted qualifiers by declaration order
	pairedPaths := make(map[string]bool)
	for index, imp := range ambiguous {
		if index < len(unaccounted) {
			pairedPaths[imp.path] = true
		}
	}

	for _, imp := range imports {
		if imp.blank || imp.dot {
			continue
		}
		if used[imp.localName] > 0 {
			continue
		}
		if pairedPaths[imp.path] {
			continue
		}
		push(types.Path{Node: &importFinding{file: file, spec: imp.spec}})
	}

	if len(unusedConstNames(file)) > 0 {
		push(p)
	}

	for _, funcDecl := range findUnusedPrivateFuncs(file, used) {
		push(types.Path{Node: &funcDeclFinding{file: file, decl: funcDecl}})
	}
}

// In the engine, `used` counts every identifier in the file outside the import
// block, including each function's own declaration name. A function is unused
// when that count never exceeds the single declaration reference.
//
// findUnusedPrivateFuncs returns unexported, non-method functions that are
// never referenced within the same file. The engine lints one file at a time,
// so a private function called from a sibling file would look unused here;
// the whole-file reference check keeps that false positive out of scope the
// same way putout handles file-scoped linting.
func findUnusedPrivateFuncs(file *ast.File, used map[string]int) []*ast.FuncDecl {
	var unused []*ast.FuncDecl
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		name := funcDecl.Name.Name
		if funcDecl.Recv != nil { // method — may satisfy an interface
			continue
		}
		if name == "init" || name == "main" { // runtime-called
			continue
		}
		if ast.IsExported(name) {
			continue
		}
		if used[name] <= 1 { // declaration only, no call or value reference
			unused = append(unused, funcDecl)
		}
	}
	return unused
}

// fixOneUnusedPrivateFunc removes a single function declaration from the file.
func fixOneUnusedPrivateFunc(file *ast.File, target *ast.FuncDecl) {
	kept := file.Decls[:0]
	for _, decl := range file.Decls {
		if decl != target {
			kept = append(kept, decl)
		}
	}
	file.Decls = kept
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
	used := countIdentUses(types.Path{Node: file}) // counts all idents outside import blocks
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
			for _, name := range spec.(*ast.ValueSpec).Names {
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
			allUnused := true
			for _, name := range spec.(*ast.ValueSpec).Names {
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
				ident := lhs.(*ast.Ident)
				if ident.Name == "_" {
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
			genDecl := declStmt.Decl.(*ast.GenDecl)
			if genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				for _, name := range spec.(*ast.ValueSpec).Names {
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
	types.Path{Node: n}.Traverse(map[string]func(types.Path){
		"*ast.Ident": func(ip types.Path) {
			reads[ip.Node.(*ast.Ident).Name]++
		},
	})
}

// findUnusedVars visits *ast.BlockStmt nodes during traversal.
func findUnusedVars(p types.Path, push func(types.Path)) {
	block := p.Node.(*ast.BlockStmt)
	if len(unusedVarNames(block)) > 0 {
		push(p)
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
				ident := lhs.(*ast.Ident)
				if ident.Name == "_" {
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
			genDecl := declStmt.Decl.(*ast.GenDecl)
			if genDecl.Tok == token.VAR {
				newSpecs := genDecl.Specs[:0]
				for _, spec := range genDecl.Specs {
					allUnused := true
					for _, name := range spec.(*ast.ValueSpec).Names {
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
