package remove_unused_variable

import "testing"

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
