// Package types contains the shared type definitions used by all plugins
// and by the engine-loader. It has no dependency on internal/engine.
package types

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"coderaiser/indra/babel"
	"coderaiser/indra/compare"
)

// Vars is the hole-bindings map from compare.
type Vars = compare.Vars

// Re-exported compare symbols for plugin dot-import convenience.
type BodySlice = compare.BodySlice
type ArgSlice = compare.ArgSlice

func Compare(node ast.Node, pattern string) bool {
	return compare.GetTemplateValues(node, pattern) != nil
}
func GetTemplateValues(node ast.Node, pattern string) map[string]ast.Node {
	return compare.GetTemplateValues(node, pattern)
}

// Re-export babel type-check functions so plugins need only dot-import types.
func IsIdent(node ast.Node) bool             { return babel.IsIdent(node) }
func IsCallExpr(node ast.Node) bool          { return babel.IsCallExpr(node) }
func IsSelector(node ast.Node) bool          { return babel.IsSelector(node) }
func IsCompositeLit(node ast.Node) bool      { return babel.IsCompositeLit(node) }
func IsArrayExpr(node ast.Node) bool         { return babel.IsArrayExpr(node) }
func IsObjectExpr(node ast.Node) bool        { return babel.IsObjectExpr(node) }
func IsFuncLit(node ast.Node) bool           { return babel.IsFuncLit(node) }
func IsBasicLit(node ast.Node) bool          { return babel.IsBasicLit(node) }
func IsStatement(node ast.Node) bool         { return babel.IsStatement(node) }
func IsFile(node ast.Node) bool              { return babel.IsFile(node) }
func IsBoolLit(node ast.Node, val bool) bool { return babel.IsBoolLit(node, val) }

// MatchFn is a guard run after pattern match. Its second argument is the Path
// of the matched statement (declaration-level matches receive a Path whose
// Node is the matched declaration). Every Matcher entry must supply a non-nil
// MatchFn. To express "no guard", omit the key from Match() entirely.
type MatchFn = func(Vars, Path) bool

// Matcher maps pattern string → optional guard. Returned by Match().
type Matcher map[string]MatchFn

// MatchWithOpts is the putout-aligned match signature: called once with the
// rule's Options, it returns a Matcher whose guards close over opts. Mirrors
// putout's match({options}) => { pattern: fn } — options are parsed once at
// plugin-load time, not threaded through the guard at each matched node.
type MatchWithOpts = func(Options) Matcher

// Replacer maps pattern string → replacement template. Returned by Replace().
type Replacer map[string]string

// Options is a per-rule config map passed to Filter functions.
// Values come from .indra.toml rule options.
type Options map[string]any

// StringSlice extracts a []string from Options by key.
// Accepts string (single value) or []string (array).
func (o Options) StringSlice(key string) []string {
	v, ok := o[key]
	if !ok {
		return nil
	}
	switch val := v.(type) {
	case string:
		return []string{val}
	case []string:
		return val
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// FilterFn is the signature for filter functions.
// Mirrors putout's filter(path, {options}) => bool.
type FilterFn = func(Vars, Path, Options) bool

// Filter is a map of pattern → FilterFn — like Matcher but receives Options.
type Filter map[string]FilterFn

// FindFn is called by the engine for each traversed node. It calls push once
// per finding with the Path the engine should pass to Fix and Report.
// A Path carries the found node and its ancestor stack, so plugins can reach
// parent selectors via Path.Find / FindParent / ParentPath.
type FindFn = func(p Path, push func(Path))

// ReportFn produces the lint message for a found path.
type ReportFn = func(p Path) string

// FixFn fixes one found path. options is per-plugin config from PluginItem.Options.
type FixFn = func(p Path, options map[string]any)

// Traverser maps key → visitor function. Returned by Traverse().
// Key formats:
//   - "*ast.File", "*ast.CallExpr", or any "*ast.*" type name:
//     visitor receives every node of that type.
//   - any other string: treated as a compare pattern (e.g. "t.Equal(__a, __b)");
//     visitor receives every *ast.ExprStmt whose AST matches the pattern.
//
// The engine merges all Traverser visitors from all plugins into one
// astutil.Apply pass. The same key formats are accepted by path.Traverse.
type Traverser map[string]FindFn

// Path is a node together with its ancestor stack, matching Babel's path API.
// Stack holds ancestors root-first, excluding Node itself; it is populated by
// the engine's PreorderStack walk and is engine-internal. Plugins reach
// ancestors only through Find / FindParent / ParentPath.
// Cursor is set by the engine during astutil.Apply and is engine-internal.
// Plugins call path.Replace / Delete / InsertBefore / InsertAfter, never
// path.Cursor directly.
// state is non-nil only for paths handed to a path.Traverse visitor; it carries
// per-call early-exit flags set by Stop and Skip. Engine-constructed paths have
// a nil state, so Stop and Skip are no-ops on them.
type Path struct {
	Node   ast.Node
	Stack  []ast.Node // ancestors root-first, excluding Node; engine-internal
	Cursor *astutil.Cursor
	state  *traverseState // non-nil only inside path.Traverse callbacks
}

// traverseState holds early-exit state for one path.Traverse call.
type traverseState struct {
	stopped bool
	skip    bool
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
// immediate parent. A visitor may call childPath.Stop() to halt the whole
// walk, or childPath.Skip() to skip only the current node's children.
// Each path.Traverse call has its own early-exit state, so stopping or
// skipping within one call never affects a later, independent Traverse call.
func (p Path) Traverse(visitors map[string]func(Path)) {
	// Split up front: type-keyed ("*ast.X") vs pattern-keyed (everything else).
	typeVisitors := make(map[string]func(Path), len(visitors))
	type patEntry struct {
		pattern string
		fn      func(Path)
	}
	var patVisitors []patEntry
	for key, fn := range visitors {
		if strings.HasPrefix(key, "*ast.") {
			typeVisitors[key] = fn
		} else {
			patVisitors = append(patVisitors, patEntry{key, fn})
		}
	}

	state := &traverseState{}
	astutil.Apply(p.Node, func(c *astutil.Cursor) bool {
		if state.stopped {
			return false
		}
		state.skip = false
		n := c.Node()
		child := Path{
			Node:   n,
			Stack:  append(append([]ast.Node{}, p.Stack...), p.Node),
			Cursor: c,
			state:  state,
		}

		// Type-keyed dispatch.
		key := fmt.Sprintf("%T", n)
		if fn, ok := typeVisitors[key]; ok {
			fn(child)
			if state.stopped {
				return false
			}
			if state.skip {
				return false
			}
		}

		// Pattern-keyed dispatch — only *ast.ExprStmt nodes can match,
		// because parsePattern wraps patterns in a func body and returns
		// Body.List[0], which is always *ast.ExprStmt.
		if _, ok := n.(*ast.ExprStmt); ok {
			for _, e := range patVisitors {
				if compare.GetTemplateValues(n, e.pattern) != nil {
					e.fn(child)
					if state.stopped {
						return false
					}
					if state.skip {
						return false
					}
				}
			}
		}

		return true
	}, nil)
}

// Stop signals path.Traverse to skip any remaining nodes in this walk,
// matching Babel's path.stop(). It is a no-op on engine-constructed paths
// that were not created inside a path.Traverse callback.
func (p Path) Stop() {
	if p.state != nil {
		p.state.stopped = true
	}
}

// Skip signals path.Traverse to skip the current node's children only;
// siblings and their descendants continue normally. It matches Babel's
// visitor returning false. It is a no-op on engine-constructed paths that
// were not created inside a path.Traverse callback.
func (p Path) Skip() {
	if p.state != nil {
		p.state.skip = true
	}
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

// PrevSibling returns the sibling path immediately before p in its parent
// slice, and true. Returns a zero Path and false when p has no prev sibling
// or its parent is not a stmt-list container. Matches Babel's
// path.getPrevSibling().
func (p Path) PrevSibling() (Path, bool) {
	if len(p.Stack) == 0 {
		return Path{}, false
	}
	parent := p.Stack[len(p.Stack)-1]
	switch par := parent.(type) {
	case *ast.BlockStmt:
		for i, stmt := range par.List {
			if stmt == p.Node && i > 0 {
				return Path{Node: par.List[i-1], Stack: p.Stack}, true
			}
		}
	case *ast.File:
		for i, decl := range par.Decls {
			if decl == p.Node && i > 0 {
				return Path{Node: par.Decls[i-1], Stack: p.Stack}, true
			}
		}
	}
	return Path{}, false
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

// Lint is the engine contract that test helpers depend on.
// Each linter implements its own and passes it to CreateTest.
// Mirrors how flatlint passes its own lint() to @putout/test.
type Lint func(src []byte, fix bool, plugins []any) (LintResult, error)

// LintResult is the minimal result shape shared across linters.
type LintResult struct {
	Out    []byte
	Places []Place
}

// Rule is a named plugin entry inside a group's Rules() return.
type Rule struct {
	Name   string
	Plugin any // carries Report/Match/Replace/Traverse/Fix methods
}
