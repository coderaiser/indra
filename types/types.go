// Package types contains the shared type definitions used by all plugins
// and by the engine-loader. It has no dependency on internal/engine.
package types

import (
	"go/ast"
	"go/token"

	"coderaiser/indra/compare"
)

// Vars is the hole-bindings map from compare.
type Vars = compare.Vars

// MatchFn is an optional guard run after pattern match.
type MatchFn = func(Vars) bool

// Matcher maps pattern string → optional guard. Returned by Match().
type Matcher map[string]MatchFn

// Replacer maps pattern string → replacement template. Returned by Replace().
type Replacer map[string]string

// VisitFn visits a matched AST node and returns findings.
type VisitFn = func(node ast.Node, vars Vars) []Place

// Traverser maps AST node type key → visitor. Returned by Traverse().
// Keys: "*ast.File", "*ast.BlockStmt"
type Traverser map[string]VisitFn

// Place is a single lint finding.
type Place struct {
	Rule    string
	Message string
	Pos     token.Position
}

// PluginEntry is a value in a Nested map.
// Enabled=true means on by default; false means off by default.
type PluginEntry struct {
	Path    string
	Enabled bool
}

// Off marks a plugin disabled by default in a Nested map.
// User can re-enable in config: "group/rule" = "on".
func Off(path string) PluginEntry {
	return PluginEntry{Path: path, Enabled: false}
}

// Nested groups sub-plugins. Values are package path strings (enabled by default)
// or PluginEntry (to set a non-default state).
type Nested map[string]any // string | PluginEntry
