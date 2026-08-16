package operator

import (
	"go/ast"
	"go/token"
	"strconv"

	"coderaiser/indra/compare"
)

// Compare reports whether node matches pattern.
func Compare(node ast.Node, pattern string) bool {
	return compare.GetTemplateValues(node, pattern) != nil
}

// GetTemplateValues matches node against pattern.
func GetTemplateValues(node ast.Node, pattern string) map[string]ast.Node {
	return compare.GetTemplateValues(node, pattern)
}

// Re-exported sentinel types from compare for plugin convenience.
type BodySlice = compare.BodySlice
type ArgSlice = compare.ArgSlice

// Extract returns the string value of node:
//   - *ast.Ident        → Name
//   - *ast.BasicLit     → Value (unquoted for strings)
//   - *ast.SelectorExpr → "X.Sel"
//
// It panics on unsupported types (same as putout's extract).
func Extract(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.BasicLit:
		if n.Kind == token.STRING {
			if s, err := strconv.Unquote(n.Value); err == nil {
				return s
			}
		}
		return n.Value
	case *ast.SelectorExpr:
		return Extract(n.X) + "." + n.Sel.Name
	default:
		panic("operator: Extract: unsupported node type")
	}
}

// IsSimple reports whether node is a literal, identifier, or selector
// expression — i.e. has no sub-expressions that could be extracted.
func IsSimple(node ast.Node) bool {
	switch node.(type) {
	case *ast.BasicLit, *ast.Ident, *ast.SelectorExpr:
		return true
	}
	return false
}

// HasImport reports whether file imports the package at importPath.
func HasImport(file *ast.File, importPath string) bool {
	if file == nil {
		return false
	}
	for _, imp := range file.Imports {
		if imp.Path != nil && imp.Path.Value == strconv.Quote(importPath) {
			return true
		}
	}
	return false
}

// GetImportAlias returns the local name for importPath in file, or "" when the
// import is absent. It returns "." for dot-imports and the alias for aliased
// imports; a plain import (no alias) returns the empty string.
func GetImportAlias(file *ast.File, importPath string) string {
	for _, imp := range file.Imports {
		if imp.Path != nil && imp.Path.Value == strconv.Quote(importPath) {
			if imp.Name == nil {
				return ""
			}
			return imp.Name.Name
		}
	}
	return ""
}

// FileFromVars extracts the *ast.File injected by the runner as "$file", or nil
// when it is not present.
func FileFromVars(vars compare.Vars) *ast.File {
	f, _ := vars["$file"].(*ast.File)
	return f
}

// BlockFromVars extracts the *ast.BlockStmt injected by the runner as "$block",
// or nil when it is not present.
func BlockFromVars(vars compare.Vars) *ast.BlockStmt {
	b, _ := vars["$block"].(*ast.BlockStmt)
	return b
}
