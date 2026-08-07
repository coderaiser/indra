package convert_no_error_to_not_ok

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

// TestFindNoErrorCallsNonFile covers the non-*ast.File early return.
func TestFindNoErrorCallsNonFile(t *testing.T) {
	Test(t, "findNoErrorCalls: non-file node is a no-op", func(t *T) {
		pushed := false
		findNoErrorCalls(types.Path{Node: ast.NewIdent("x")}, func(types.Path) { pushed = true })
		t.Ok(!pushed)
		t.End()
	})
}

// TestFixNonFile covers the non-*ast.File early return.
func TestFixNonFile(t *testing.T) {
	Test(t, "Fix: non-file node is a no-op", func(t *T) {
		Fix(types.Path{Node: ast.NewIdent("x")}, nil)
		t.Pass("no panic")
		t.End()
	})
}
