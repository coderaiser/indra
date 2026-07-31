package engine

import (
	"go/ast"
	"go/token"
	"strings"
	"testing"
)

func reportOnlyPlugin() Plugin {
	return Plugin{
		Name:   "report",
		Report: func() string { return "message" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"t.Equal(__a, __b)": nil}
		},
	}
}

func replacePlugin() Plugin {
	return Plugin{
		Name:   "replace",
		Report: func() string { return "message" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"t.Equal(__a, __b)": nil}
		},
		Replace: func() map[string]string {
			return map[string]string{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}
		},
	}
}

func traversePlugin() Plugin {
	return Plugin{
		Name:   "traverse",
		Report: func() string { return "file issue" },
		Traverse: func() map[string]TraverseVisitor {
			return map[string]TraverseVisitor{
				"*ast.File": func(node ast.Node, vars Vars) []Place {
					return []Place{{Message: "file issue", Pos: token.Position{Line: 1}}}
				},
			}
		},
	}
}

const equalSrc = `package p

import "testing"

func TestFoo(t *testing.T) {
	t.Equal(a, b)
}
`

func TestReportOnly(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []Plugin{reportOnlyPlugin()}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if places[0].Rule != "report" || places[0].Message != "message" {
		t.Fatalf("unexpected place: %+v", places[0])
	}
	if string(out) != equalSrc {
		t.Fatal("report-only must leave src unchanged")
	}
}

func TestReplacePlugin(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []Plugin{replacePlugin()}, true)
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
	out, places, err := Indra([]byte(equalSrc), []Plugin{traversePlugin()}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if places[0].Rule != "traverse" || places[0].Message != "file issue" {
		t.Fatalf("unexpected place: %+v", places[0])
	}
	if string(out) != equalSrc {
		t.Fatal("traverse plugin must leave src unchanged")
	}
}

func TestFixFalseKeepsReplacement(t *testing.T) {
	out, places, err := Indra([]byte(equalSrc), []Plugin{replacePlugin()}, false)
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
	p := Plugin{
		Name:   "multi",
		Report: func() string { return "msg" },
		Match: func() map[string]MatchFn {
			return map[string]MatchFn{"makeSlices(__x)": nil}
		},
		Replace: func() map[string]string {
			return map[string]string{"makeSlices(__x)": "x := __x\ny := __x"}
		},
	}
	src := `package p

func f() {
	makeSlices(v)
}
`
	out, _, err := Indra([]byte(src), []Plugin{p}, true)
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, "x := v") || !strings.Contains(got, "y := v") {
		t.Fatalf("expected both statements:\n%s", got)
	}
}

func TestMultiplePluginsOrder(t *testing.T) {
	p2 := reportOnlyPlugin()
	p2.Name = "second"
	_, places, err := Indra([]byte(equalSrc), []Plugin{reportOnlyPlugin(), p2}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 2 {
		t.Fatalf("expected 2 places, got %d", len(places))
	}
	if places[0].Rule != "report" || places[1].Rule != "second" {
		t.Fatalf("plugins ran out of order: %+v", places)
	}
}

func TestParseError(t *testing.T) {
	src := []byte("package p\nfunc (\n")
	out, places, err := Indra(src, []Plugin{reportOnlyPlugin()}, false)
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
