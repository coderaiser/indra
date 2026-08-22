package apply_exports_order_test

import (
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/internal/plugin_indra/apply_exports_order"
	. "coderaiser/indra/internal/test"
	. "coderaiser/indra/types"
)

var Test = CreateTest("apply-exports-order", apply_exports_order.Plugin{})

func TestApplyExportsOrder(t *testing.T) {
	Test(t, "apply-exports-order: report: out-of-order", func(t *T) {
		t.Report("out-of-order", "Apply exports order")
		t.End()
	})

	Test(t, "apply-exports-order: transform: out-of-order", func(t *T) {
		t.Transform("out-of-order")
		t.End()
	})

	Test(t, "apply-exports-order: no report: in-order", func(t *T) {
		t.NoReport("in-order")
		t.End()
	})

	Test(t, "apply-exports-order: no transform: in-order", func(t *T) {
		t.NoTransform("in-order")
		t.End()
	})

	Test(t, "apply-exports-order: no report: no-shape", func(t *T) {
		t.NoReport("no-shape")
		t.End()
	})

	Test(t, "apply-exports-order: no transform: no-shape", func(t *T) {
		t.NoTransform("no-shape")
		t.End()
	})
}

// TestFixNoShapeKeepsDecls checks that Fix leaves a file untouched when its
// exported functions do not implement any known plugin shape. The runner only
// calls Fix on pushed paths, so this branch needs a direct call.
func TestFixNoShapeKeepsDecls(t *testing.T) {
	src := "package p\n\nvar keep = 1\n\nfunc helper() {}\n"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	apply_exports_order.Fix(Path{Node: file}, nil)
	if len(file.Decls) != 2 {
		t.Fatalf("expected 2 decls, got %d", len(file.Decls))
	}
}
