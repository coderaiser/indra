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

// Path is a node together with its ancestor stack, matching Babel's path API.
// Stack holds ancestors root-first, excluding Node itself; it is populated by
// the engine's PreorderStack walk and is engine-internal. Plugins reach
// ancestors only through Find / FindParent / ParentPath.
type Path struct {
	Node  ast.Node
	Stack []ast.Node // ancestors root-first, excluding Node; engine-internal
}

// Find walks up from this path (inclusive) and returns the first Path where fn
// returns true. Returns a zero Path and false when no ancestor matches.
// Matches Babel's path.find(fn).
func (p Path) Find(fn func(Path) bool) (Path, bool) {
	if fn(p) {
		return p, true
	}
	return p.FindParent(fn)
}

// FindParent walks up from p's parent (exclusive of self) and returns the
// nearest ancestor Path where fn returns true. Matches Babel's
// path.findParent(fn).
func (p Path) FindParent(fn func(Path) bool) (Path, bool) {
	for i := len(p.Stack) - 1; i >= 0; i-- {
		ancestor := Path{Node: p.Stack[i], Stack: p.Stack[:i]}
		if fn(ancestor) {
			return ancestor, true
		}
	}
	return Path{}, false
}

// ParentPath returns the immediate parent Path, or a zero Path and false when
// this path has no parent. Matches Babel's path.parentPath.
func (p Path) ParentPath() (Path, bool) {
	if len(p.Stack) == 0 {
		return Path{}, false
	}
	parent := p.Stack[len(p.Stack)-1]
	return Path{Node: parent, Stack: p.Stack[:len(p.Stack)-1]}, true
}

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
