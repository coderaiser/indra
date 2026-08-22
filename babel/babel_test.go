package babel

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	. "github.com/coderaiser/go-tape"
)

// fakeNode is a custom ast.Node carrying a slice of non-Node values, used to
// exercise GetAll's element-type guard.
type fakeNode struct {
	ast.Node
	Tags []string
}

// parseSrc parses src into an *ast.File for traverse tests.
func parseSrc(t *testing.T, src string) *ast.File {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

func TestTraverse(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tdo()\n}\n"
	file := parseSrc(t, src)

	files := 0
	funcs := 0
	calls := 0
	hits := 0
	Traverse(file, map[string]func(Path){
		"File":     func(p Path) { files++ },
		"FuncDecl": func(p Path) { funcs++ },
		"CallExpr": func(p Path) { calls++ },
	})
	Traverse(file, map[string]func(Path){
		"BasicLit": func(p Path) { hits++ },
	})

	Test(t, "babel: Traverse: visits File nodes", func(t *T) {
		t.Equal(files, 1)
		t.End()
	})

	Test(t, "babel: Traverse: visits FuncDecl nodes", func(t *T) {
		t.Equal(funcs, 1)
		t.End()
	})

	Test(t, "babel: Traverse: visits CallExpr nodes", func(t *T) {
		t.Equal(calls, 1)
		t.End()
	})

	Test(t, "babel: Traverse: ignores unvisited type names", func(t *T) {
		t.NotOk(hits)

		t.End()
	})
}

func TestPathGet(t *testing.T) {
	call := &ast.CallExpr{
		Fun:  ast.NewIdent("do"),
		Args: []ast.Expr{ast.NewIdent("a")},
	}

	Test(t, "babel: Get: returns the child at a named field", func(t *T) {
		fun := Path{Node: call}.Get("Fun")
		t.Ok(!fun.IsEmpty() && IsIdent(fun.Node))
		t.End()
	})

	Test(t, "babel: Get: empty path for missing field", func(t *T) {
		t.Ok(Path{Node: call}.Get("Nope").IsEmpty())
		t.End()
	})

	Test(t, "babel: Get: empty path when node is nil", func(t *T) {
		t.Ok(Path{}.Get("Fun").IsEmpty())
		t.End()
	})

	Test(t, "babel: Get: empty path when node is a typed nil", func(t *T) {
		t.Ok(Path{Node: (*ast.Ident)(nil)}.Get("Name").IsEmpty())
		t.End()
	})

	Test(t, "babel: Get: empty path when field is nil", func(t *T) {
		sel := &ast.SelectorExpr{X: ast.NewIdent("x")}
		t.Ok(Path{Node: sel}.Get("Sel").IsEmpty())
		t.End()
	})
}

func TestPathGetAll(t *testing.T) {
	call := &ast.CallExpr{
		Fun:  ast.NewIdent("do"),
		Args: []ast.Expr{ast.NewIdent("a"), ast.NewIdent("b")},
	}

	Test(t, "babel: GetAll: returns children of a slice field", func(t *T) {
		args := Path{Node: call}.GetAll("Args")
		t.Ok(len(args) == 2 && IsIdent(args[0].Node))
		t.End()
	})

	Test(t, "babel: GetAll: nil for missing and non-slice fields", func(t *T) {
		p := Path{Node: call}
		t.Ok(p.GetAll("Nope") == nil && p.GetAll("Fun") == nil)
		t.End()
	})

	Test(t, "babel: GetAll: nil when node is empty", func(t *T) {
		t.NotOk(Path{}.GetAll("Args"))

		t.End()
	})

	Test(t, "babel: GetAll: nil when node is a typed nil", func(t *T) {
		t.NotOk(Path{Node: (*ast.Ident)(nil)}.GetAll("Args"))

		t.End()
	})

	Test(t, "babel: GetAll: skips elements that are not ast.Node", func(t *T) {
		fake := &fakeNode{Tags: []string{"a", "b"}}
		tags := Path{Node: fake}.GetAll("Tags")
		t.NotOk(len(tags))

		t.End()
	})
}

func TestPathType(t *testing.T) {
	Test(t, "babel: Type: returns the node type name", func(t *T) {
		result := Path{Node: &ast.CallExpr{}}.Type()
		t.Equal(result, "CallExpr")

		t.End()
	})

	Test(t, "babel: Type: empty for a zero path", func(t *T) {
		t.NotOk(Path{}.Type())

		t.End()
	})
}

func TestPathIsEmpty(t *testing.T) {
	Test(t, "babel: IsEmpty: true for zero path", func(t *T) {
		t.Ok(Path{}.IsEmpty())
		t.End()
	})

	Test(t, "babel: IsEmpty: false when node set", func(t *T) {
		t.NotOk(Path{Node: ast.NewIdent("x")}.IsEmpty())
		t.End()
	})
}

func TestIsIdent(t *testing.T) {
	Test(t, "babel: IsIdent: identifier", func(t *T) {
		t.Ok(IsIdent(ast.NewIdent("x")))
		t.End()
	})

	Test(t, "babel: IsIdent: non-identifier", func(t *T) {
		t.NotOk(IsIdent(&ast.BasicLit{Kind: token.STRING, Value: `"x"`}))
		t.End()
	})
}

func TestIsCallExpr(t *testing.T) {
	Test(t, "babel: IsCallExpr: call expression", func(t *T) {
		t.Ok(IsCallExpr(&ast.CallExpr{Fun: ast.NewIdent("f")}))
		t.End()
	})

	Test(t, "babel: IsCallExpr: non-call", func(t *T) {
		t.NotOk(IsCallExpr(ast.NewIdent("f")))
		t.End()
	})
}

func TestIsSelector(t *testing.T) {
	Test(t, "babel: IsSelector: selector expression", func(t *T) {
		sel := &ast.SelectorExpr{X: ast.NewIdent("x"), Sel: ast.NewIdent("Y")}
		t.Ok(IsSelector(sel))
		t.End()
	})

	Test(t, "babel: IsSelector: non-selector", func(t *T) {
		t.NotOk(IsSelector(ast.NewIdent("x")))
		t.End()
	})
}

func TestIsCompositeLit(t *testing.T) {
	Test(t, "babel: IsCompositeLit: composite literal", func(t *T) {
		t.Ok(IsCompositeLit(&ast.CompositeLit{Type: ast.NewIdent("T")}))
		t.End()
	})

	Test(t, "babel: IsCompositeLit: non-composite-literal", func(t *T) {
		t.NotOk(IsCompositeLit(ast.NewIdent("T")))
		t.End()
	})
}

func TestIsArrayExpr(t *testing.T) {
	Test(t, "babel: IsArrayExpr: slice literal", func(t *T) {
		t.Ok(IsArrayExpr(&ast.CompositeLit{Type: &ast.ArrayType{}}))
		t.End()
	})

	Test(t, "babel: IsArrayExpr: struct literal is not a slice", func(t *T) {
		t.NotOk(IsArrayExpr(&ast.CompositeLit{Type: ast.NewIdent("T")}))
		t.End()
	})

	Test(t, "babel: IsArrayExpr: type-less literal is not a slice", func(t *T) {
		t.NotOk(IsArrayExpr(&ast.CompositeLit{}))
		t.End()
	})

	Test(t, "babel: IsArrayExpr: non-composite-literal", func(t *T) {
		t.NotOk(IsArrayExpr(ast.NewIdent("T")))
		t.End()
	})
}

func TestIsObjectExpr(t *testing.T) {
	Test(t, "babel: IsObjectExpr: struct literal", func(t *T) {
		t.Ok(IsObjectExpr(&ast.CompositeLit{Type: ast.NewIdent("T")}))
		t.End()
	})

	Test(t, "babel: IsObjectExpr: qualified struct literal", func(t *T) {
		lit := &ast.CompositeLit{
			Type: &ast.SelectorExpr{X: ast.NewIdent("pkg"), Sel: ast.NewIdent("T")},
		}
		t.Ok(IsObjectExpr(lit))
		t.End()
	})

	Test(t, "babel: IsObjectExpr: slice literal is not an object", func(t *T) {
		t.NotOk(IsObjectExpr(&ast.CompositeLit{Type: &ast.ArrayType{}}))
		t.End()
	})

	Test(t, "babel: IsObjectExpr: type-less literal is not an object", func(t *T) {
		t.NotOk(IsObjectExpr(&ast.CompositeLit{}))
		t.End()
	})

	Test(t, "babel: IsObjectExpr: non-composite-literal", func(t *T) {
		t.NotOk(IsObjectExpr(ast.NewIdent("T")))
		t.End()
	})
}

func TestIsFuncLit(t *testing.T) {
	Test(t, "babel: IsFuncLit: function literal", func(t *T) {
		t.Ok(IsFuncLit(&ast.FuncLit{Type: &ast.FuncType{}}))
		t.End()
	})

	Test(t, "babel: IsFuncLit: non-func-literal", func(t *T) {
		t.NotOk(IsFuncLit(ast.NewIdent("f")))
		t.End()
	})
}

func TestIsBasicLit(t *testing.T) {
	Test(t, "babel: IsBasicLit: basic literal", func(t *T) {
		t.Ok(IsBasicLit(&ast.BasicLit{Kind: token.INT, Value: "1"}))
		t.End()
	})

	Test(t, "babel: IsBasicLit: non-basic-literal", func(t *T) {
		t.NotOk(IsBasicLit(ast.NewIdent("1")))
		t.End()
	})
}

func TestIsStatement(t *testing.T) {
	Test(t, "babel: IsStatement: statement", func(t *T) {
		t.Ok(IsStatement(&ast.EmptyStmt{}))
		t.End()
	})

	Test(t, "babel: IsStatement: expression is not a statement", func(t *T) {
		t.NotOk(IsStatement(&ast.CallExpr{Fun: ast.NewIdent("f")}))
		t.End()
	})
}

func TestIsFile(t *testing.T) {
	file := parseSrc(t, "package p\n")

	Test(t, "babel: IsFile: file", func(t *T) {
		t.Ok(IsFile(file))
		t.End()
	})

	Test(t, "babel: IsFile: non-file", func(t *T) {
		t.NotOk(IsFile(&ast.EmptyStmt{}))
		t.End()
	})
}

func TestIsBoolLit(t *testing.T) {
	Test(t, "IsBoolLit: true literal matches true", func(t *T) {
		t.Ok(IsBoolLit(ast.NewIdent("true"), true))
		t.End()
	})

	Test(t, "IsBoolLit: true literal rejects false", func(t *T) {
		t.NotOk(IsBoolLit(ast.NewIdent("true"), false))
		t.End()
	})

	Test(t, "IsBoolLit: false literal matches false", func(t *T) {
		t.Ok(IsBoolLit(ast.NewIdent("false"), false))
		t.End()
	})

	Test(t, "IsBoolLit: false literal rejects true", func(t *T) {
		t.NotOk(IsBoolLit(ast.NewIdent("false"), true))
		t.End()
	})

	Test(t, "IsBoolLit: non-identifier is not a bool literal", func(t *T) {
		t.NotOk(IsBoolLit(&ast.BasicLit{Kind: token.STRING, Value: `"true"`}, true))
		t.End()
	})

	Test(t, "IsBoolLit: other identifier is not a bool literal", func(t *T) {
		t.NotOk(IsBoolLit(ast.NewIdent("t"), true))
		t.End()
	})
}
