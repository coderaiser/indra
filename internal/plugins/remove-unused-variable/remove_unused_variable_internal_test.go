package remove_unused_variable

import (
	"go/ast"
	"go/token"
	"testing"

	"coderaiser/indra/types"
)

// TestSelfReportMessage verifies the self wrapper forwards the report message,
// covering the method the engine-loader reflects on.
func TestSelfReportMessage(t *testing.T) {
	if got := Self.Report(); got != "remove unused variable" {
		t.Fatalf("unexpected self report message: %q", got)
	}
}

// TestSelfTraverseCovered verifies the self wrapper exposes the block visitor.
func TestSelfTraverseCovered(t *testing.T) {
	tr := Self.Traverse()
	if _, ok := tr["*ast.BlockStmt"]; !ok {
		t.Fatal("expected self.Traverse to expose *ast.BlockStmt visitor")
	}
}

// TestTopLevelReportMessage covers the top-level Report func directly.
func TestTopLevelReportMessage(t *testing.T) {
	if got := Report(); got != "remove unused variable" {
		t.Fatalf("unexpected report message: %q", got)
	}
}

// TestSelfFixRemovesUnusedVariable covers self.Fix, the method the runner
// reflects on, by removing an unused variable from a hand-built block.
func TestSelfFixRemovesUnusedVariable(t *testing.T) {
	block := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Tok: token.DEFINE,
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "1"}},
			},
		},
	}
	Self.Fix(block, []types.Place{{Message: "remove unused variable: x"}})
	if len(block.List) != 0 {
		t.Fatalf("expected unused variable statement removed, got %d stmts", len(block.List))
	}
}
