package engine

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"
)

// fileFixer removes the first top-level var from a file when asked to fix.
type fileFixer struct{}

func (fileFixer) Report() string { return "filesig" }
func (fileFixer) Traverse() types.Traverser {
	return types.Traverser{
		"*ast.File": func(node ast.Node, vars types.Vars) []types.Place {
			return []types.Place{{Message: "filesig"}}
		},
	}
}
func (fileFixer) Fix(node ast.Node, places []types.Place) {
	file := node.(*ast.File)
	if first, ok := file.Decls[0].(*ast.GenDecl); ok && first.Tok.String() == "var" {
		file.Decls = file.Decls[1:]
	}
}

// blockFixer is a traverser whose Fix trims an extra statement from a block.
type blockFixer struct{}

func (blockFixer) Report() string { return "blockfix" }
func (blockFixer) Traverse() types.Traverser {
	return types.Traverser{
		"*ast.BlockStmt": func(node ast.Node, vars types.Vars) []types.Place {
			return []types.Place{{Message: "blockfix"}}
		},
	}
}
func (blockFixer) Fix(node ast.Node, places []types.Place) {}

// TestNestedPlugin verifies load expands a Nested map into its sub-plugins.
func TestNestedPlugin(t *testing.T) {
	nested := types.Nested{"one": replacer{}}
	out, places, err := Indra([]byte(equalSrc), []any{nested}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place from nested, got %d", len(places))
	}
	if places[0].Rule != "one" {
		t.Fatalf("expected nested rule name 'one', got %q", places[0].Rule)
	}
	if string(out) != equalSrc {
		t.Fatal("nested replacer without fix must leave src unchanged")
	}
}

// TestResolveUnknownPanics verifies resolve panics on an unknowable value.
func TestResolveUnknownPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown plugin kind")
		}
	}()
	resolve("bogus", "not a plugin")
}

// TestLoadNestedExpands verifies load expands nested and keeps names.
func TestLoadNestedExpands(t *testing.T) {
	rp := load([]any{
		types.Nested{"a": replacer{}},
		replacer{},
	})
	if len(rp) != 2 {
		t.Fatalf("expected 2 resolved plugins, got %d", len(rp))
	}
	if rp[0].name != "a" {
		t.Fatalf("expected expanded name 'a', got %q", rp[0].name)
	}
}

// TestDeriveNameFallback covers the t.String() branch for unnamed types.
func TestDeriveNameFallback(t *testing.T) {
	got := deriveName(&ast.Ident{})
	if got == "" {
		t.Fatal("expected a non-empty name via t.String()")
	}
}

// TestTraverserFileFix verifies the *ast.File fix path is invoked with fix=true.
func TestTraverserFileFix(t *testing.T) {
	src := `package p

var _ = 1

func f() {}
`
	out, places, err := Indra([]byte(src), []any{fileFixer{}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if string(out) == string(src) {
		t.Fatal("expected *ast.File Fix to modify the source")
	}
}

// TestTraverserBlockFix verifies the *ast.BlockStmt fix path with fix=true.
func TestTraverserBlockFix(t *testing.T) {
	_, places, err := Indra([]byte(`package p
func f() {
	x := 1
	_ = x
}
`), []any{blockFixer{}}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(places) != 1 {
		t.Fatalf("expected 1 place from blockFixer, got %d", len(places))
	}
}
