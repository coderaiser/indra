package add_t_end

import (
	"go/ast"
	"testing"

	"coderaiser/indra/types"
)

// TestMissingEndBadBody covers the defensive !ok branch of missingEnd: when
// the matched __body hole is not a BodySlice (which can't happen via the engine
// because Match only fires once the pattern has bound __body to one), the guard
// must reject the node rather than panicking.
func TestMissingEndBadBody(t *testing.T) {
	if missingEnd(types.Vars{"__body": &ast.Ident{Name: "x"}}, types.Path{}) {
		t.Fatal("expected missingEnd to return false for a non-BodySlice __body")
	}
}
