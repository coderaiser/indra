package convert_no_error_to_not_ok

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	. "github.com/coderaiser/go-tape"
)

func TestFixWithoutGoTape(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "", "package p\nfunc f() { t.NoError(err) }\n", 0)
	if err != nil {
		t.Fatal(err)
	}
	// No go-tape import: Fix must return early and change nothing.
	Fix(file, nil)

	Test(t, "fix: leaves NoError when go-tape not imported", func(t *T) {
		t.Ok(hasSelectorName(file, "NoError"))
		t.End()
	})
}

// hasSelectorName reports whether any selector in file uses Sel name.
func hasSelectorName(file *ast.File, name string) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if ok && sel.Sel != nil && sel.Sel.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}
