// Package types contains the shared type definitions used by all plugins
// and by the engine-loader. It has no dependency on internal/engine.
package types

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"

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
// per finding with the Path the engine should pass to Fix and Report.
// A Path carries the found node and its ancestor stack, so plugins can reach
// parent selectors via Path.Find / FindParent / ParentPath.
type FindFn = func(p Path, push func(Path))

// ReportFn produces the lint message for a found path.
type ReportFn = func(p Path) string

// FixFn fixes one found path. options is per-plugin config from PluginItem.Options.
type FixFn = func(p Path, options map[string]any)

// Traverser maps AST node type key → finder. Returned by Traverse().
// Keys: "*ast.File", "*ast.BlockStmt"
type Traverser map[string]FindFn

// Path is a node together with its ancestor stack, matching Babel's path API.
// Stack holds ancestors root-first, excluding Node itself; it is populated by
// the engine's PreorderStack walk and is engine-internal. Plugins reach
// ancestors only through Find / FindParent / ParentPath.
// Cursor is set by the engine during astutil.Apply and is engine-internal.
// Plugins call path.Replace / Delete / InsertBefore / InsertAfter, never
// path.Cursor directly.
type Path struct {
	Node   ast.Node
	Stack  []ast.Node // ancestors root-first, excluding Node; engine-internal
	Cursor *astutil.Cursor
}

// Replace delegates to Cursor.Replace when Cursor is non-nil.
func (p Path) Replace(n ast.Node) {
	if p.Cursor != nil {
		p.Cursor.Replace(n)
	}
}

// Delete delegates to Cursor.Delete when Cursor is non-nil.
func (p Path) Delete() {
	if p.Cursor != nil {
		p.Cursor.Delete()
	}
}

// InsertBefore delegates to Cursor.InsertBefore when Cursor is non-nil.
func (p Path) InsertBefore(n ast.Node) {
	if p.Cursor != nil {
		p.Cursor.InsertBefore(n)
	}
}

// InsertAfter delegates to Cursor.InsertAfter when Cursor is non-nil.
func (p Path) InsertAfter(n ast.Node) {
	if p.Cursor != nil {
		p.Cursor.InsertAfter(n)
	}
}

// Traverse walks the sub-tree rooted at p.Node in pre-order, routing each
// node to the matching visitor by its type key ("*ast.ReturnStmt" etc).
// Visitor keys use the same format as Traverser.
// The visitor receives a child Path whose Stack includes p.Node as the
// immediate parent. There is no early-exit mechanism in this first version;
// if early exit is needed, use a closed-over bool and return immediately.
func (p Path) Traverse(visitors map[string]func(Path)) {
	astutil.Apply(p.Node, func(c *astutil.Cursor) bool {
		key := fmt.Sprintf("%T", c.Node())
		if fn, ok := visitors[key]; ok {
			child := Path{
				Node:   c.Node(),
				Stack:  append(append([]ast.Node{}, p.Stack...), p.Node),
				Cursor: c,
			}
			fn(child)
		}
		return true
	}, nil)
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
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Place is a single lint finding. JSON tags mirror putout's formatter shape
// (rule/message/position) so json-lines output is drop-in compatible.
type Place struct {
	Rule     string   `json:"rule"`
	Message  string   `json:"message"`
	Position Position `json:"position"`
}

// Rule is a named plugin entry inside a group's Rules() return.
type Rule struct {
	Name   string
	Plugin any // carries Report/Match/Replace/Traverse/Fix methods
}
