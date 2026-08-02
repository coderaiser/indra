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

// Nested is a group plugin: maps rule name → sub-plugin Self value.
type Nested map[string]any
