// Package apply_compare migrates plugins from the low-level compare API to the
// higher-level operator API. It rewrites
//
//	compare.GetTemplateValues(x, p) != nil   →   Compare(x, p)
//
// and retargets the compare import to the operator dot-import so the file
// keeps compiling with a single plugin import.
package apply_compare

import (
	"go/ast"
	"go/token"
	"strings"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "use Compare instead of GetTemplateValues != nil" }

// Fix rewrites the found binary into a bare Compare(...) call and retargets the
// compare import to the operator dot-import.
func Fix(p Path, _ map[string]any) {
	binary := p.Node.(*ast.BinaryExpr)
	call := binary.X.(*ast.CallExpr)
	call.Fun = ast.NewIdent("Compare")
	p.Replace(call)
	file, _ := p.FindParent(func(parent Path) bool {
		_, isFile := parent.Node.(*ast.File)
		return isFile
	})
	retargetImport(file.Node.(*ast.File))
}

func Traverse() Traverser {
	return Traverser{"*ast.File": findUselessComparisons}
}

// findUselessComparisons pushes a Path for every
// compare.GetTemplateValues(x, p) != nil binary expression inside a plugin
// file. Plugin files are recognized by their Plugin struct marker, so engine
// internals that legitimately use compare (e.g. engine_runner) are skipped.
func findUselessComparisons(p Path, push func(Path)) {
	file := p.Node.(*ast.File)
	if !hasPluginStruct(file) {
		return
	}
	p.Traverse(map[string]func(Path){
		"*ast.BinaryExpr": func(binaryPath Path) {
			if isCompareNil(binaryPath.Node.(*ast.BinaryExpr)) {
				push(binaryPath)
			}
		},
	})
}

// hasPluginStruct reports whether any top-level declaration is a type named
// Plugin — the conventional marker of an indra/plugin file.
func hasPluginStruct(file *ast.File) bool {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == "Plugin" {
				return true
			}
		}
	}
	return false
}

// isCompareNil reports whether binary is a `compare.GetTemplateValues(...) != nil`
// expression: a not-equal comparison whose left side is a GetTemplateValues
// call on the compare selector and whose right side is nil.
func isCompareNil(binary *ast.BinaryExpr) bool {
	if binary.Op != token.NEQ {
		return false
	}
	call, ok := binary.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || id.Name != "compare" {
		return false
	}
	if sel.Sel == nil || sel.Sel.Name != "GetTemplateValues" {
		return false
	}
	nilIdent, ok := binary.Y.(*ast.Ident)
	return ok && nilIdent.Name == "nil"
}

// retargetImport turns the "coderaiser/indra/compare" import spec into the
// operator dot-import.
func retargetImport(file *ast.File) {
	for _, imp := range file.Imports {
		if strings.Trim(imp.Path.Value, `"`) != "coderaiser/indra/compare" {
			continue
		}
		imp.Name = ast.NewIdent(".")
		imp.Path.Value = `"coderaiser/indra/operator"`
		return
	}
}

// Plugin wraps the rule for the registry: a fixing AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
