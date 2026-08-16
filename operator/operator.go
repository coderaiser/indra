package operator

import (
	"go/ast"

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
