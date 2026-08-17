package add_t_end

import (
	"testing"

	"coderaiser/indra/types"
)

// TestMissingEndBadBody covers the defensive !ok branch of missingEnd: when
// the matched __body hole is not a BodySlice (which can't happen via the engine
// because Match only fires once the pattern has bound __body to one), the guard
// must reject the node rather than panicking.
func TestMissingEndBadBody(t *testing.T) {
	if missingEnd(types.Vars{"__body": "not-a-body"}, types.Path{}) {
		t.Fatal("expected missingEnd to return false for a non-BodySlice __body")
	}
}

// TestTapeImportedNoFile covers tapeImported's fallback when the matched path
// has no enclosing *ast.File on its stack.
func TestTapeImportedNoFile(t *testing.T) {
	if tapeImported(types.Vars{}, types.Path{Stack: []ast.Node{&ast.Ident{}}}) {
		t.Fatal("expected tapeImported to return false without an *ast.File")
	}
}
