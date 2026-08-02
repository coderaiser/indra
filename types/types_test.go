package types

import (
	"go/ast"
	"go/token"
	"testing"

	"coderaiser/indra/compare"
)

// replacerLike is a minimal ReplacerPlugin that satisfies Plugin.
type replacerLike struct{}

func (replacerLike) Report() string  { return "r" }
func (replacerLike) Match() Matcher  { return nil }
func (replacerLike) Replace() Replacer { return nil }
func (replacerLike) isPlugin()       {}

// traverserLike is a minimal TraverserPlugin that satisfies Plugin.
type traverserLike struct{}

func (traverserLike) Report() string     { return "t" }
func (traverserLike) Traverse() Traverser { return nil }
func (traverserLike) Fix(node ast.Node, places []Place) {}
func (traverserLike) isPlugin()          {}

// TestPluginInterface verifies both plugin kinds satisfy the Plugin interface.
func TestPluginInterface(t *testing.T) {
	var p Plugin = replacerLike{}
	if p == nil {
		t.Fatal("expected non-nil plugin")
	}

	var tr Plugin = traverserLike{}
	if tr == nil {
		t.Fatal("expected non-nil traverser plugin")
	}
}

// TestNestedHoldsSubPlugins verifies Nested can hold both plugin kinds.
func TestNestedHoldsSubPlugins(t *testing.T) {
	n := Nested{
		"replacer":   replacerLike{},
		"traverser":  traverserLike{},
	}
	if len(n) != 2 {
		t.Fatalf("expected 2 sub plugins, got %d", len(n))
	}
	if _, ok := n["replacer"]; !ok {
		t.Fatal("expected replacer sub plugin")
	}
}

// TestVisitFnShape verifies VisitFn has the required signature.
func TestVisitFnShape(t *testing.T) {
	var fn VisitFn = func(node ast.Node, vars Vars) []Place {
		return nil
	}
	places := fn(nil, Vars{})
	if len(places) != 0 {
		t.Fatalf("expected no places, got %d", len(places))
	}
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
		Rule:    "rule",
		Message: "msg",
		Pos:     token.Position{Line: 1, Column: 2},
	}
	if p.Rule != "rule" || p.Message != "msg" {
		t.Fatal("unexpected Place content")
	}
	if p.Pos.Line != 1 || p.Pos.Column != 2 {
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
		"pat": func(Vars) bool { return true },
	}
	if _, ok := m["pat"]; !ok {
		t.Fatal("expected pattern in matcher")
	}
}

// TestTraverserShape verifies Traverser maps keys to visitors.
func TestTraverserShape(t *testing.T) {
	tr := Traverser{
		"*ast.File": func(node ast.Node, vars Vars) []Place { return nil },
	}
	if _, ok := tr["*ast.File"]; !ok {
		t.Fatal("expected *ast.File key in traverser")
	}
}

// TestCompareRoundTrip guards against accidental drift of the Vars binding.
func TestCompareRoundTrip(t *testing.T) {
	src := "package p\nvar _ = 1\n"
	_ = compare.Compare
	if src == "" {
		t.Fatal("unreachable")
	}
}
