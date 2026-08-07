package types_test

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"
	. "github.com/coderaiser/go-tape"
)

func makeStack(nodes ...ast.Node) types.Path {
	if len(nodes) == 0 {
		return types.Path{}
	}
	return types.Path{
		Node:  nodes[len(nodes)-1],
		Stack: nodes[:len(nodes)-1],
	}
}

func ident(name string) *ast.Ident { return ast.NewIdent(name) }

func TestFind(t *testing.T) {
	a, b, c := ident("a"), ident("b"), ident("c")
	path := types.Path{Node: c, Stack: []ast.Node{a, b}}

	Test(t, "Path.Find: matches self", func(t *T) {
		_, ok := path.Find(func(p types.Path) bool { return p.Node == c })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.Find: matches ancestor", func(t *T) {
		_, ok := path.Find(func(p types.Path) bool { return p.Node == a })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.Find: returns false when nothing matches", func(t *T) {
		_, ok := path.Find(func(p types.Path) bool { return false })
		t.NotOk(ok)

		t.End()
	})
}

func TestFindParent(t *testing.T) {
	a, b, c := ident("a"), ident("b"), ident("c")
	path := types.Path{Node: c, Stack: []ast.Node{a, b}}

	Test(t, "Path.FindParent: finds immediate parent", func(t *T) {
		_, ok := path.FindParent(func(p types.Path) bool { return p.Node == b })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.FindParent: finds grandparent", func(t *T) {
		_, ok := path.FindParent(func(p types.Path) bool { return p.Node == a })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.FindParent: does not match self", func(t *T) {
		_, ok := path.FindParent(func(p types.Path) bool { return p.Node == c })
		t.NotOk(ok)

		t.End()
	})

	Test(t, "Path.FindParent: returns false on empty stack", func(t *T) {
		root := types.Path{Node: a, Stack: nil}
		_, ok := root.FindParent(func(p types.Path) bool { return true })
		t.NotOk(ok)

		t.End()
	})

	Test(t, "Path.FindParent: ancestor stack is trimmed correctly", func(t *T) {
		found, _ := path.FindParent(func(p types.Path) bool { return p.Node == b })
		result := len(found.Stack)
		t.Equal(result, 1)

		t.End()
	})
}

func TestParentPath(t *testing.T) {
	a, b := ident("a"), ident("b")
	path := types.Path{Node: b, Stack: []ast.Node{a}}

	Test(t, "Path.ParentPath: returns parent", func(t *T) {
		_, ok := path.ParentPath()
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.ParentPath: returns false at root", func(t *T) {
		_, ok := types.Path{Node: a, Stack: nil}.ParentPath()
		t.NotOk(ok)

		t.End()
	})

	Test(t, "Path.ParentPath: parent has empty stack at root", func(t *T) {
		parent, _ := path.ParentPath()
		result := len(parent.Stack)
		t.Equal(result, 0)

		t.End()
	})
}
