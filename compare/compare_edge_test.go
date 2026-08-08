package compare

import (
	"testing"
)

func TestMatchNodeArgsNotCall(t *testing.T) {
	vars := make(Vars)
	// __args pattern against a non-CallExpr real falls through and mismatches.
	if matchNode(parseStmt(t, "t.Equal(__args)"), parseStmt(t, "foo"), vars) {
		t.Fatal("__args against non-call should not match")
	}
}

func TestMatchNodeArgsLiteral(t *testing.T) {
	// a single non-ident argument disables the __args fast path.
	if GetTemplateValues(parseStmt(t, "t.Equal(1)"), "t.Equal(1)") == nil {
		t.Fatal("literal single-arg call should match")
	}
}

func TestMatchNodeBodySelectorFun(t *testing.T) {
	// the sentinel call fun being a selector (not ident) skips the capture.
	node := parseStmt(t, "{\na.b()\n}")
	if GetTemplateValues(node, "{\na.b()\n}") == nil {
		t.Fatal("selector call inside block should match")
	}
}
