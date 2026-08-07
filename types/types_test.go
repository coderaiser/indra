package types

import (
	"go/ast"
	"testing"
)

// replacerLike is a minimal plugin exposing the replacer method set.
type replacerLike struct{}

func (replacerLike) Report() string    { return "r" }
func (replacerLike) Match() Matcher    { return nil }
func (replacerLike) Replace() Replacer { return nil }

// traverserLike is a minimal plugin exposing the traverser method set.
type traverserLike struct{}

func (traverserLike) Report() string                         { return "t" }
func (traverserLike) Traverse() Traverser                    { return nil }
func (traverserLike) Fix(node ast.Node, opts map[string]any) {}

// TestRuleDefaults verifies the Rule struct carries a name and plugin.
func TestRuleDefaults(t *testing.T) {
	r := Rule{Name: "remove-skip", Plugin: replacerLike{}}
	if r.Name != "remove-skip" {
		t.Fatalf("unexpected Rule defaults: %+v", r)
	}
	if _, ok := r.Plugin.(replacerLike); !ok {
		t.Fatalf("expected replacer plugin, got %T", r.Plugin)
	}
}

// TestFindFnShape verifies FindFn is assignable from a node+push callback.
func TestFindFnShape(t *testing.T) {
	var fn = func(node ast.Node, push func(ast.Node)) {
		push(node)
	}
	var got ast.Node
	fn(&ast.File{}, func(n ast.Node) { got = n })
	if got == nil {
		t.Fatal("expected push to be called with the node")
	}
}

// TestReportFnShape verifies ReportFn has the required signature.
func TestReportFnShape(t *testing.T) {
	var fn = func(node ast.Node) string { return "msg" }
	if fn(nil) != "msg" {
		t.Fatal("unexpected report result")
	}
}

// TestFixFnShape verifies FixFn has the required signature.
func TestFixFnShape(t *testing.T) {
	var fn = func(node ast.Node, options map[string]any) {}
	fn(nil, nil)
}

// TestVarsAlias verifies Vars aliases compare.Vars.
func TestVarsAlias(t *testing.T) {
	v := make(Vars)
	v["__a"] = &ast.Ident{Name: "x"}
	_, ok := v["__a"]
	if !ok {
		t.Fatal("expected key __a")
	}
}

// TestPlaceFields verifies Place carries rule, message and position.
func TestPlaceFields(t *testing.T) {
	p := Place{
		Rule:     "rule",
		Message:  "msg",
		Position: Position{Line: 1, Column: 2},
	}
	if p.Rule != "rule" || p.Message != "msg" {
		t.Fatal("unexpected Place content")
	}
	if p.Position.Line != 1 || p.Position.Column != 2 {
		t.Fatal("unexpected Place position")
	}
}

// TestReplacerShape verifies Replacer is a string map.
func TestReplacerShape(t *testing.T) {
	r := Replacer{"a": "b"}
	if r["a"] != "b" {
		t.Fatal("unexpected replacer value")
	}
}

// TestMatcherShape verifies Matcher is a map of pattern to guard.
func TestMatcherShape(t *testing.T) {
	m := Matcher{
		"pat": func(Vars, *ast.BlockStmt) bool { return true },
	}
	if _, ok := m["pat"]; !ok {
		t.Fatal("expected pattern in matcher")
	}
}

// TestTraverserShape verifies Traverser maps keys to finders.
func TestTraverserShape(t *testing.T) {
	tr := Traverser{
		"*ast.File": func(p Path, push func(Path)) {},
	}
		if _, ok := tr["*ast.File"]; !ok {
		t.Fatal("expected *ast.File key in traverser")
	}
}

// TestPathFindReturnsSelf verifies Find returns the path itself when fn matches.
func TestPathFindReturnsSelf(t *testing.T) {
	self := &ast.Ident{Name: "leaf"}
	p := Path{Node: self}
	found, ok := p.Find(func(anc Path) bool { return anc.Node == self })
	if !ok || found.Node != self {
		t.Fatalf("expected Find to return self, got ok=%v node=%v", ok, found.Node)
	}
}

// TestPathFindUp verifies Find walks up to a matching ancestor.
func TestPathFindUp(t *testing.T) {
	root := &ast.Ident{Name: "root"}
	mid := &ast.Ident{Name: "mid"}
	p := Path{Node: &ast.Ident{Name: "leaf"}, Stack: []ast.Node{root, mid}}
	found, ok := p.Find(func(anc Path) bool { return anc.Node == mid })
	if !ok || found.Node != mid {
		t.Fatalf("expected Find to return mid, got ok=%v node=%v", ok, found.Node)
	}
}

// TestPathFindUpNotFound verifies Find returns false when no ancestor matches.
func TestPathFindUpNotFound(t *testing.T) {
	p := Path{Node: &ast.Ident{Name: "leaf"}, Stack: []ast.Node{&ast.Ident{Name: "root"}}}
	_, ok := p.Find(func(Path) bool { return false })
	if ok {
		t.Fatal("expected Find to return false when no ancestor matches")
	}
}

// TestPathFindParent verifies FindParent returns the nearest ancestor, not self.
func TestPathFindParent(t *testing.T) {
	leaf := &ast.Ident{Name: "leaf"}
	root := &ast.Ident{Name: "root"}
	p := Path{Node: leaf, Stack: []ast.Node{root}}
	found, ok := p.FindParent(func(Path) bool { return true })
	if !ok || found.Node != root {
		t.Fatalf("expected nearest parent root, got ok=%v node=%v", ok, found.Node)
	}
}

// TestPathParentPath verifies ParentPath returns the immediate parent.
func TestPathParentPath(t *testing.T) {
	root := &ast.Ident{Name: "root"}
	p := Path{Node: &ast.Ident{Name: "leaf"}, Stack: []ast.Node{root}}
	parent, ok := p.ParentPath()
	if !ok || parent.Node != root {
		t.Fatalf("expected immediate parent root, got ok=%v node=%v", ok, parent.Node)
	}
}

// TestPathParentPathEmpty verifies ParentPath returns false for a root path.
func TestPathParentPathEmpty(t *testing.T) {
	p := Path{Node: &ast.Ident{Name: "root"}}
	_, ok := p.ParentPath()
	if ok {
		t.Fatal("expected ParentPath to return false for empty stack")
	}
}
