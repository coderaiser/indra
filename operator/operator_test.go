package operator_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	. "coderaiser/indra/operator"

	. "github.com/coderaiser/go-tape"
)

func parseStmt(src string) ast.Node {
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", "package p\nfunc _(){"+src+"}", 0)
	return file.Decls[0].(*ast.FuncDecl).Body.List[0]
}

func TestOperator(t *testing.T) {
	Test(t, "operator: Compare returns true on match", func(t *T) {
		t.Ok(Compare(parseStmt("t.End()"), "__.End()"))
		t.End()
	})

	Test(t, "operator: Compare returns false on no match", func(t *T) {
		t.NotOk(Compare(parseStmt("t.Begin()"), "__.End()"))

		t.End()
	})

	Test(t, "operator: GetTemplateValues returns vars on match", func(t *T) {
		t.Ok(GetTemplateValues(parseStmt("t.Ok(err)"), "__a.Ok(__b)"))

		t.End()
	})

	Test(t, "operator: GetTemplateValues returns nil on no match", func(t *T) {
		t.NotOk(GetTemplateValues(parseStmt("t.End()"), "__a.Ok(__b)"))

		t.End()
	})
}
