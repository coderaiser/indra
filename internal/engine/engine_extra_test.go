package engine

import (
	"go/ast"
	"strings"
	"testing"
)

func TestNonMatchingStatement(t *testing.T) {
	src := `package p

func f() {
	foo()
	t.Equal(a, b)
}
`
	_, places, err := Indra([]byte(src), []Plugin{reportOnlyPlugin()}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}

func TestGuardRejects(t *testing.T) {
	p := Plugin{
		Name:   "guarded",
		Report: func() string { return "msg" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"t.Equal(__a, __b)": func(v Vars) bool { return false }}
		},
	}
	_, places, err := Indra([]byte(equalSrc), []Plugin{p}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 0 {
		t.Fatalf("expected 0 places, got %d", len(places))
	}
}

func TestGuardAccepts(t *testing.T) {
	p := Plugin{
		Name:   "guarded",
		Report: func() string { return "msg" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"t.Equal(__a, __b)": func(v Vars) bool { return true }}
		},
	}
	_, places, err := Indra([]byte(equalSrc), []Plugin{p}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}

func TestInvalidTemplateSkipsRewrite(t *testing.T) {
	p := Plugin{
		Name:   "bad",
		Report: func() string { return "msg" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"call(__a)": nil}
		},
		Replace: func() map[string]string {
			return map[string]string{"call(__a)": "this is ( broken"}
		},
	}
	out, places, err := Indra([]byte("package p\nfunc f() {\n\tcall(x)\n}\n"), []Plugin{p}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if strings.Contains(string(out), "broken") {
		t.Fatal("invalid template must not be applied")
	}
}

func TestApplyRewritesDescending(t *testing.T) {
	src := `package p

func f() {
	t.Equal(a, b)
	t.Equal(c, d)
}
`
	out, _, err := Indra([]byte(src), []Plugin{replacePlugin()}, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "t.DeepEqual(a, b)") || !strings.Contains(got, "t.DeepEqual(c, d)") {
		t.Fatalf("expected both rewrites:\n%s", got)
	}
}

func TestRenderArgsSlice(t *testing.T) {
	p := Plugin{
		Name:   "wrap",
		Report: func() string { return "msg" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"g(__args)": nil}
		},
		Replace: func() map[string]string {
			return map[string]string{"g(__args)": "f(__args)"}
		},
	}
	out, _, err := Indra([]byte("package p\nfunc f() {\n\tg(a, b)\n}\n"), []Plugin{p}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "f(a, b)") {
		t.Fatalf("expected f(a, b):\n%s", out)
	}
}

func TestRenderBodySlice(t *testing.T) {
	p := Plugin{
		Name:   "body",
		Report: func() string { return "msg" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"g(func() { __body })": nil}
		},
		Replace: func() map[string]string {
			return map[string]string{"g(func() { __body })": "h(func() {\n__body\n})"}
		},
	}
	out, _, err := Indra([]byte("package p\nfunc f() {\n\tg(func() { x() })\n}\n"), []Plugin{p}, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "h(func()") || !strings.Contains(got, "x()") {
		t.Fatalf("expected body rewrite:\n%s", got)
	}
}

func TestSubstituteHelpers(t *testing.T) {
	if substituteAndParse("this is ( broken", Vars{}) != nil {
		t.Fatal("invalid template should return nil")
	}
	if got := substitute("a __nope b", Vars{}); got != "a __nope b" {
		t.Fatalf("unknown hole should be left alone, got %q", got)
	}
	if got := substitute("x __y z", Vars{"__y": ast.NewIdent("w")}); got != "x w z" {
		t.Fatalf("expected substitution, got %q", got)
	}
	if printNode(nil) != "" {
		t.Fatal("printNode(nil) should be empty")
	}
}

func TestTraverseBlockStmt(t *testing.T) {
	p := Plugin{
		Name:   "block-visitor",
		Report: func() string { return "block issue" },
		Traverse: func() map[string]TraverseVisitor {
			return map[string]TraverseVisitor{
				"*ast.BlockStmt": func(node ast.Node, vars Vars) []Place {
					return []Place{{Message: "block issue"}}
				},
			}
		},
	}
	src := `package p
func f() {
	x := 1
	_ = x
}
`
	_, places, err := Indra([]byte(src), []Plugin{p}, false)
	if err != nil {
		t.Fatal(err)
	}
	// f() has one block, so expect 1 place
	if len(places) != 1 {
		t.Fatalf("expected 1 place from block visitor, got %d", len(places))
	}
	if places[0].Rule != "block-visitor" {
		t.Fatalf("unexpected rule: %s", places[0].Rule)
	}
}

func TestStripPositionsNilNode(t *testing.T) {
	// ast.Inspect calls the visitor with nil to signal end-of-subtree.
	// Pass a node whose child is nil to exercise the nil branch in stripPositions.
	node := &ast.Ident{Name: "x"}
	// stripPositions must not panic on a node that ast.Inspect will call with nil
	stripPositions(node)
}
