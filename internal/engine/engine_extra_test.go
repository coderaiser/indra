package engine

import (
	"go/ast"
	"strings"
	"testing"

	"coderaiser/indra/types"
)

// second is a duplicate of reportOnly to test ordering.
type second struct{}

func (second) Report() string { return "message" }
func (second) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": nil}
}
func (second) Replace() types.Replacer { return nil }

// multiPlacer replaces one statement with two.
type multiPlacer struct{}

func (multiPlacer) Report() string { return "msg" }
func (multiPlacer) Match() types.Matcher {
	return types.Matcher{"makeSlices(__x)": nil}
}
func (multiPlacer) Replace() types.Replacer {
	return types.Replacer{"makeSlices(__x)": "x := __x\ny := __x"}
}

// guardedFalse has a guard that rejects every match.
type guardedFalse struct{}

func (guardedFalse) Report() string { return "msg" }
func (guardedFalse) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return false }}
}
func (guardedFalse) Replace() types.Replacer { return nil }

// guardedTrue has a guard that accepts every match.
type guardedTrue struct{}

func (guardedTrue) Report() string { return "msg" }
func (guardedTrue) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}
}
func (guardedTrue) Replace() types.Replacer { return nil }

// badTemplate has an invalid replacement template.
type badTemplate struct{}

func (badTemplate) Report() string { return "msg" }
func (badTemplate) Match() types.Matcher {
	return types.Matcher{"call(__a)": nil}
}
func (badTemplate) Replace() types.Replacer {
	return types.Replacer{"call(__a)": "this is ( broken"}
}

// wrap moves an arg slice into a call.
type wrap struct{}

func (wrap) Report() string { return "msg" }
func (wrap) Match() types.Matcher {
	return types.Matcher{"g(__args)": nil}
}
func (wrap) Replace() types.Replacer {
	return types.Replacer{"g(__args)": "f(__args)"}
}

// body rewraps a function literal body.
type body struct{}

func (body) Report() string { return "msg" }
func (body) Match() types.Matcher {
	return types.Matcher{"g(func() { __body })": nil}
}
func (body) Replace() types.Replacer {
	return types.Replacer{"g(func() { __body })": "h(func() {\n__body\n})"}
}

// blockVisitor reports one place for every block.
type blockVisitor struct{}

func (blockVisitor) Report() string { return "block issue" }
func (blockVisitor) Traverse() types.Traverser {
	return types.Traverser{
		"*ast.BlockStmt": func(node ast.Node, vars types.Vars) []types.Place {
			return []types.Place{{Message: "block issue"}}
		},
	}
}
func (blockVisitor) Fix(node ast.Node, places []types.Place) {}

func TestNonMatchingStatement(t *testing.T) {
	src := `package p

func f() {
	foo()
	t.Equal(a, b)
}
`
	_, places, err := Indra([]byte(src), []any{reportOnly{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}

func TestGuardRejects(t *testing.T) {
	_, places, err := Indra([]byte(equalSrc), []any{guardedFalse{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 0 {
		t.Fatalf("expected 0 places, got %d", len(places))
	}
}

func TestGuardAccepts(t *testing.T) {
	_, places, err := Indra([]byte(equalSrc), []any{guardedTrue{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
}

func TestInvalidTemplateSkipsRewrite(t *testing.T) {
	out, places, err := Indra([]byte("package p\nfunc f() {\n\tcall(x)\n}\n"), []any{badTemplate{}}, true)
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
	out, _, err := Indra([]byte(src), []any{replacer{}}, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "t.DeepEqual(a, b)") || !strings.Contains(got, "t.DeepEqual(c, d)") {
		t.Fatalf("expected both rewrites:\n%s", got)
	}
}

func TestRenderArgsSlice(t *testing.T) {
	out, _, err := Indra([]byte("package p\nfunc f() {\n\tg(a, b)\n}\n"), []any{wrap{}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "f(a, b)") {
		t.Fatalf("expected f(a, b):\n%s", out)
	}
}

func TestRenderBodySlice(t *testing.T) {
	out, _, err := Indra([]byte("package p\nfunc f() {\n\tg(func() { x() })\n}\n"), []any{body{}}, true)
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
	src := `package p
func f() {
	x := 1
	_ = x
}
`
	_, places, err := Indra([]byte(src), []any{blockVisitor{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place from block visitor, got %d", len(places))
	}
	if places[0].Rule != "blockVisitor" {
		t.Fatalf("unexpected rule: %s", places[0].Rule)
	}
}

func TestStripPositionsNilNode(t *testing.T) {
	node := &ast.Ident{Name: "x"}
	stripPositions(node)
}
