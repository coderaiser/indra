package remove_useless_match

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/types"
)

// matchDecl builds a func Match() Matcher { return <lit> } FuncDecl.
// Pass nil lit to build a body that returns a non-CompositeLit.
func matchDecl(lit ast.Expr) *ast.FuncDecl {
	var results []ast.Expr
	if lit != nil {
		results = []ast.Expr{lit}
	} else {
		results = []ast.Expr{ast.NewIdent("Matcher")}
	}
	return &ast.FuncDecl{
		Name: ast.NewIdent("Match"),
		Type: &ast.FuncType{
			Results: &ast.FieldList{List: []*ast.Field{
				{Type: ast.NewIdent("Matcher")},
			}},
		},
		Body: &ast.BlockStmt{List: []ast.Stmt{
			&ast.ReturnStmt{Results: results},
		}},
	}
}

func TestMatcherLitEmpty(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	lit := matcherLit(fn)
	if lit == nil {
		t.Fatal("expected empty composite literal")
	}
}

func TestMatcherLitNonComposite(t *testing.T) {
	fn := matchDecl(ast.NewIdent("Matcher"))
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for non-composite-literal return")
	}
}

func TestMatcherLitWithRecv(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Recv = &ast.FieldList{}
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for method receiver")
	}
}

func TestMatcherLitWithNamedReturn(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Type.Results.List[0].Names = []*ast.Ident{ast.NewIdent("m")}
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for named return")
	}
}

func TestMatcherLitWrongResultType(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Type.Results.List[0].Type = ast.NewIdent("Replacer")
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for non-Matcher result type")
	}
}

func TestMatcherLitNoBody(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Body = nil
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for nil body")
	}
}

func TestMatcherLitBodyWithoutReturn(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Body.List = []ast.Stmt{&ast.ExprStmt{X: ast.NewIdent("x")}}
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for non-return body")
	}
}

func TestMatcherLitTypeNil(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Type = nil
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for nil Type")
	}
}

func TestMatcherLitResultsNil(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Type.Results = nil
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for nil Results")
	}
}

func TestMatcherLitTwoResults(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	fn.Type.Results.List = append(fn.Type.Results.List, &ast.Field{Type: ast.NewIdent("error")})
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for two result fields")
	}
}

func TestMatcherLitMultiReturnStmt(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	ret := fn.Body.List[0].(*ast.ReturnStmt)
	ret.Results = append(ret.Results, ast.NewIdent("Matcher"))
	if matcherLit(fn) != nil {
		t.Fatal("expected nil for multi-value return")
	}
}

func TestIsUselessMatchEmpty(t *testing.T) {
	fn := matchDecl(&ast.CompositeLit{})
	if !isUselessMatch(fn) {
		t.Fatal("expected empty Matcher to be useless")
	}
}

func TestIsUselessMatchNilGuard(t *testing.T) {
	lit := &ast.CompositeLit{Elts: []ast.Expr{
		&ast.KeyValueExpr{Key: ast.NewIdent(`"p"`), Value: ast.NewIdent("nil")},
	}}
	if !isUselessMatch(matchDecl(lit)) {
		t.Fatal("expected all-nil Matcher to be useless")
	}
}

func TestIsUselessMatchNonKeyValue(t *testing.T) {
	lit := &ast.CompositeLit{Elts: []ast.Expr{ast.NewIdent("x")}}
	if isUselessMatch(matchDecl(lit)) {
		t.Fatal("expected non-KeyValueExpr element to be useful")
	}
}

func TestIsUselessMatchRealGuard(t *testing.T) {
	lit := &ast.CompositeLit{Elts: []ast.Expr{
		&ast.KeyValueExpr{Key: ast.NewIdent(`"p"`), Value: ast.NewIdent("f")},
	}}
	if isUselessMatch(matchDecl(lit)) {
		t.Fatal("expected non-nil guard to be useful")
	}
}

func TestIsUselessMatchNotMatchShape(t *testing.T) {
	fn := &ast.FuncDecl{Name: ast.NewIdent("Other")}
	if isUselessMatch(fn) {
		t.Fatal("expected non-Match-shaped decl to be useful")
	}
}

func TestIsUselessMatchMixedGuards(t *testing.T) {
	lit := &ast.CompositeLit{Elts: []ast.Expr{
		&ast.KeyValueExpr{Key: ast.NewIdent(`"a"`), Value: ast.NewIdent("nil")},
		&ast.KeyValueExpr{Key: ast.NewIdent(`"b"`), Value: ast.NewIdent("f")},
	}}
	if isUselessMatch(matchDecl(lit)) {
		t.Fatal("expected mixed guard Matcher to be useful")
	}
}

// TestFindUselessMatchNonFile covers the non-*ast.File early return.
func TestFindUselessMatchNonFile(t *testing.T) {
	Test(t, "findUselessMatch: non-file node is a no-op", func(t *T) {
		pushed := false
		findUselessMatch(types.Path{Node: ast.NewIdent("x")}, func(types.Path) { pushed = true })
		t.Ok(!pushed)
		t.End()
	})
}

// TestFindUselessMatchNoMatch covers the loop where no Match decl is useless.
func TestFindUselessMatchNoMatch(t *testing.T) {
	Test(t, "findUselessMatch: file without useless Match pushes nothing", func(t *T) {
		pushed := false
		file := parseFile(t, "package p\nfunc Match() Matcher { return Matcher{}}\n")
		findUselessMatch(types.Path{Node: file}, func(types.Path) { pushed = true })
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

// parseFile parses src into an *ast.File for direct helper tests.
func parseFile(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "fixture.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}
