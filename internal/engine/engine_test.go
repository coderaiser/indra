package engine

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"coderaiser/indra/types"
)

// reportOnly is a Self-shaped replacer plugin with no Replace method.
type reportOnly struct{}

func (reportOnly) Report() string { return "message" }
func (reportOnly) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": nil}
}
func (reportOnly) Replace() types.Replacer { return nil }

// replacer is a Self-shaped replacer plugin with a Replace template.
type replacer struct{}

func (replacer) Report() string { return "message" }
func (replacer) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": nil}
}
func (replacer) Replace() types.Replacer {
	return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}
}

// traverser is a Self-shaped traverser plugin reporting on the whole file.
type traverser struct{}

func (traverser) Report() string { return "file issue" }
func (traverser) Traverse() types.Traverser {
	return types.Traverser{
		"*ast.File": func(node ast.Node, vars types.Vars) []types.Place {
			return []types.Place{{Message: "file issue", Pos: token.Position{Line: 1}}}
		},
	}
}
func (traverser) Fix(node ast.Node, places []types.Place) {}

const equalSrc = `package p

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(a, b)
}
`

func TestReportOnly(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []any{reportOnly{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if places[0].Rule != "reportOnly" || places[0].Message != "message" {
		t.Fatalf("unexpected place: %+v", places[0])
	}
	if string(out) != equalSrc {
		t.Fatal("report-only must leave src unchanged")
	}
}

func TestReplacePlugin(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []any{replacer{}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if !strings.Contains(string(out), "t.DeepEqual(a, b)") {
		t.Fatalf("expected DeepEqual in output:\n%s", out)
	}
}

func TestTraversePlugin(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []any{traverser{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if places[0].Rule != "traverser" || places[0].Message != "file issue" {
		t.Fatalf("unexpected place: %+v", places[0])
	}
	if string(out) != equalSrc {
		t.Fatal("traverse plugin must leave src unchanged")
	}
}

func TestFixFalseKeepsReplacement(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []any{replacer{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if string(out) != equalSrc {
		t.Fatal("fix=false must leave src unchanged even with Replace")
	}
}

func TestMultiStmtReplace(t *testing.T) {
	p := multiPlacer{}
	src := `package p

func f() {
	makeSlices(v)
}
`
	out, _, err := Indra([]byte(src), []any{p}, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "x := v") || !strings.Contains(got, "y := v") {
		t.Fatalf("expected both statements:\n%s", got)
	}
}

func TestMultiplePluginsOrder(t *testing.T) {
	_, places, err := Indra([]byte(equalSrc), []any{reportOnly{}, second{}}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 2 {
		t.Fatalf("expected 2 places, got %d", len(places))
	}
	if places[0].Rule != "reportOnly" || places[1].Rule != "second" {
		t.Fatalf("plugins ran out of order: %+v", places)
	}
}

func TestParseError(t *testing.T) {
	src := []byte("package p\nfunc (\n")
	out, places, err := Indra(src, []any{reportOnly{}}, false)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if places != nil {
		t.Fatal("expected no places on parse error")
	}
	if string(out) != string(src) {
		t.Fatal("src should be returned untouched on parse error")
	}
}
