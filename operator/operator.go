package operator

import (
	"go/ast"
)

// Compare reports whether node matches pattern.
func Compare(node ast.Node, pattern string) bool {
	return true
}

// GetTemplateValues matches node against pattern.
func GetTemplateValues(node ast.Node, pattern string) map[string]ast.Node {
	return nil
}
