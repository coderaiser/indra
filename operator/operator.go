// Package operator provides the public plugin API for pattern matching.
// Plugins dot-import this package instead of reaching into compare directly.
//
//	import . "coderaiser/indra/operator"
//
//	if Compare(path.Node, "t.End()") { ... }
package operator

import (
	"go/ast"

	"coderaiser/indra/compare"
)

// Compare reports whether node matches pattern.
// Boolean convenience over GetTemplateValues.
func Compare(node ast.Node, pattern string) bool {
	return compare.GetTemplateValues(node, pattern) != nil
}

// GetTemplateValues matches node against pattern and returns bound holes.
// Returns nil when there is no match. Mirrors putout's getTemplateValues.
func GetTemplateValues(node ast.Node, pattern string) compare.Vars {
	return compare.GetTemplateValues(node, pattern)
}

// BodySlice is the bound node type for the __body hole. Re-exported from
// compare so plugins keep a single import.
type BodySlice = compare.BodySlice

// ArgSlice is the bound node type for the __args hole. Re-exported from
// compare so plugins keep a single import.
type ArgSlice = compare.ArgSlice
