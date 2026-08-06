// Package types contains the shared type definitions used by all plugins
// and by the engine-loader. It has no dependency on internal/engine.
package types

import (
	"go/ast"

	"coderaiser/indra/compare"
)

// Vars is the hole-bindings map from compare.
type Vars = compare.Vars

// MatchFn is a guard run after pattern match. Its second argument is the
// *ast.BlockStmt containing the matched statement (nil for declaration-level
// matches). Every Matcher entry must supply a non-nil MatchFn. To express
// "no guard", omit the key from Match() entirely.
type MatchFn = func(Vars, *ast.BlockStmt) bool

// Matcher maps pattern string → optional guard. Returned by Match().
type Matcher map[string]MatchFn

// Replacer maps pattern string → replacement template. Returned by Replace().
type Replacer map[string]string

// FindFn is called by the engine for each traversed node. It calls push once
// per finding with the node the engine should pass to Fix and Report.
type FindFn = func(node ast.Node, push func(ast.Node))

// ReportFn produces the lint message for a found node.
type ReportFn = func(node ast.Node) string

// FixFn fixes one found node. options is per-plugin config from PluginItem.Options.
type FixFn = func(node ast.Node, options map[string]any)

// Traverser maps AST node type key → finder. Returned by Traverse().
// Keys: "*ast.File", "*ast.BlockStmt"
type Traverser map[string]FindFn

// Position is a source location — line and column only, matching putout's shape.
type Position struct {
	Line   int
	Column int
}

// Place is a single lint finding.
type Place struct {
	Rule     string
	Message  string
	Position Position
}

// Rule is a named plugin entry inside a group's Rules() return.
type Rule struct {
	Name   string
	Plugin any // carries Report/Match/Replace/Traverse/Fix methods
}
