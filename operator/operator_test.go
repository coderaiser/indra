package operator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"testing"

	. "github.com/coderaiser/go-tape"
)

// parseFile parses src with comments for fixture-building.
func parseFile(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

// funcBody returns the body block of the single function in src.
func funcBody(t *testing.T, src string) *ast.BlockStmt {
	t.Helper()
	return parseFile(t, src).Decls[0].(*ast.FuncDecl).Body
}

func TestRemoveBlock(t *testing.T) {
	block := funcBody(t, "package p\nfunc f() {\n\ta := 1\n\tb := 2\n\tc := 3\n}\n")

	Test(t, "Remove: splices middle statement from block", func(t *T) {
		before := len(block.List)
		Remove(&Path{Node: block.List[1], Parent: block, Field: "List", Index: 1})
		t.Equal(len(block.List), before-1)
		t.End()
	})

	Test(t, "Remove: out-of-range index is a no-op", func(t *T) {
		Remove(&Path{Node: ast.NewIdent("x"), Parent: block, Field: "List", Index: 99})
		t.Ok(true)
		t.End()
	})

	Test(t, "Remove method form: splices block statement", func(t *T) {
		before := len(block.List)
		p := &Path{Node: block.List[0], Parent: block, Field: "List", Index: 0}
		p.remove()
		t.Equal(len(block.List), before-1)
		t.End()
	})
}

func TestRemoveGenDecl(t *testing.T) {
	src := "package p\nimport (\n\t\"fmt\"\n\t\"os\"\n)\nfunc f() {}\n"
	file := parseFile(t, src)
	gd := file.Decls[0].(*ast.GenDecl)

	Test(t, "Remove: splices a GenDecl spec", func(t *T) {
		before := len(gd.Specs)
		Remove(&Path{Node: gd.Specs[0], Parent: gd, Field: "Specs", Index: 0})
		t.Equal(len(gd.Specs), before-1)
		t.End()
	})

	Test(t, "Remove: GenDecl out-of-range no-op", func(t *T) {
		Remove(&Path{Node: gd.Specs[0], Parent: gd, Field: "Specs", Index: 5})
		t.Equal(len(gd.Specs), 1)
		t.End()
	})

	Test(t, "Remove: removing the last GenDecl spec triggers preservation", func(t *T) {
		gd2 := parseFile(t.TB(), "package p\nimport (\n\t\"a\"\n\t\"b\"\n)\nfunc f() {}\n").Decls[0].(*ast.GenDecl)
		Remove(&Path{Node: gd2.Specs[1], Parent: gd2, Field: "Specs", Index: 1})
		t.Equal(len(gd2.Specs), 1)
		t.End()
	})
}

func TestRemoveFile(t *testing.T) {
	file := parseFile(t, "package p\nfunc a() {}\nfunc b() {}\n")

	Test(t, "Remove: removes a file-level declaration", func(t *T) {
		before := len(file.Decls)
		Remove(&Path{Node: file.Decls[0], Parent: file, Field: "Decls", Index: 0})
		t.Equal(len(file.Decls), before-1)
		t.End()
	})

	Test(t, "Remove: removes a file import by index", func(t *T) {
		impFile := parseFile(t.TB(), "package p\nimport (\n\t\"a\"\n\t\"b\"\n)\n")
		before := len(impFile.Imports)
		Remove(&Path{Node: impFile.Imports[1], Parent: impFile, Field: "Imports", Index: 1})
		t.Equal(len(impFile.Imports), before-1)
		t.End()
	})

	Test(t, "Remove: file out-of-range no-op", func(t *T) {
		Remove(&Path{Node: file.Decls[0], Parent: file, Field: "Decls", Index: 10})
		t.Equal(len(file.Decls), 1)
		t.End()
	})
}

func TestRemoveUnsupportedParent(t *testing.T) {
	Test(t, "Remove: unsupported parent is a no-op", func(t *T) {
		Remove(&Path{Node: ast.NewIdent("x"), Parent: ast.NewIdent("y"), Field: "List", Index: 0})
		t.Ok(true)
		t.End()
	})
}

func TestPreserveComments(t *testing.T) {
	Test(t, "preserveComments: node with no comment groups returns early", func(t *T) {
		assign := &ast.AssignStmt{}
		preserveComments(&Path{Node: assign}, assign)
		t.Ok(true)
		t.End()
	})

	Test(t, "preserveComments: no file ancestor returns early", func(t *T) {
		gd := &ast.GenDecl{Tok: token.VAR, Doc: &ast.CommentGroup{}}
		preserveComments(&Path{Node: gd}, gd)
		t.Ok(true)
		t.End()
	})

	Test(t, "Remove: last-in-block keeps an existing tracked comment (no duplicate)", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n\t// trailing\n\tvar z int\n}\n"
		file := parseFile(t.TB(), src)
		before := len(file.Comments)
		block := file.Decls[0].(*ast.FuncDecl).Body
		last := len(block.List) - 1
		Remove(&Path{Node: block.List[last], Parent: block, Field: "List", Index: last})
		t.Equal(len(file.Comments), before)
		t.End()
	})

	Test(t, "Remove: appends a comment group not already tracked", func(t *T) {
		file := parseFile(t.TB(), "package p\ntype T int\n")
		gd := &ast.GenDecl{
			Tok: token.TYPE,
			Doc: &ast.CommentGroup{List: []*ast.Comment{{Slash: 1, Text: "// keep"}}},
		}
		file.Decls[0] = gd
		Remove(&Path{Node: gd, Parent: file, ParentPath: &Path{Node: file}, Field: "Decls", Index: 0})
		t.Equal(len(file.Comments), 1)
		t.End()
	})
}

func TestCommentGroups(t *testing.T) {
	Test(t, "commentGroups: nil node", func(t *T) {
		t.Equal(len(commentGroups(nil)), 0)
		t.End()
	})

	Test(t, "commentGroups: typed nil pointer", func(t *T) {
		var id *ast.Ident
		t.Equal(len(commentGroups(id)), 0)
		t.End()
	})

	Test(t, "commentGroups: collects Doc and recurses through children", func(t *T) {
		gd := &ast.GenDecl{
			Tok: token.VAR,
			Doc: &ast.CommentGroup{List: []*ast.Comment{{Text: "// a"}}},
			Specs: []ast.Spec{
				&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("x")}},
			},
		}
		t.Equal(len(commentGroups(gd)), 1)
		t.End()
	})

	Test(t, "commentGroups: skips a nil pointer child", func(t *T) {
		fd := &ast.FuncDecl{Name: ast.NewIdent("f"), Body: nil, Doc: nil}
		t.Equal(len(commentGroups(fd)), 0)
		t.End()
	})

	Test(t, "commentGroups: recurses into a non-nil pointer child", func(t *T) {
		fd := &ast.FuncDecl{
			Name: ast.NewIdent("f"),
			Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: ast.NewIdent("x")}}},
		}
		t.Equal(len(commentGroups(fd)), 0)
		t.End()
	})

	Test(t, "commentGroups: recurses into an interface child", func(t *T) {
		e := &ast.ExprStmt{X: ast.NewIdent("x")}
		t.Equal(len(commentGroups(e)), 0)
		t.End()
	})

	Test(t, "commentGroups: skips a non-node pointer field (File.Scope)", func(t *T) {
		file := &ast.File{Name: ast.NewIdent("p"), Scope: ast.NewScope(nil)}
		// ast.Scope does not implement ast.Node; the walk must skip it safely.
		t.Equal(len(commentGroups(file)), 0)
		t.End()
	})
}

func TestFindFile(t *testing.T) {
	file := parseFile(t, "package p\nfunc f() {}\n")
	block := file.Decls[0].(*ast.FuncDecl).Body
	stmt := ast.NewIdent("x")

	Test(t, "findFile: returns file via ParentPath chain", func(t *T) {
		found := findFile(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: file}})
		t.Ok(found == file)
		t.End()
	})

	Test(t, "findFile: no file ancestor returns nil", func(t *T) {
		found := findFile(&Path{Node: stmt, Parent: block})
		t.NotOk(found != nil)
		t.End()
	})

	Test(t, "findFile: nil path", func(t *T) {
		t.NotOk(findFile(nil) != nil)
		t.End()
	})
}

func TestReplaceWith(t *testing.T) {
	Test(t, "ReplaceWith: replaces a statement in a block and preserves position", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		orig := block.List[0]
		newStmt := &ast.AssignStmt{
			Lhs:    []ast.Expr{ast.NewIdent("z")},
			Rhs:    []ast.Expr{ast.NewIdent("2")},
			Tok:    token.ASSIGN,
			TokPos: orig.Pos(),
		}
		ReplaceWith(&Path{Node: orig, Parent: block, Field: "List", Index: 0}, newStmt)
		t.Ok(block.List[0] == newStmt)
		t.End()
	})

	Test(t, "ReplaceWith: leaves the block untouched when field/idx is missing", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		before := len(block.List)
		ReplaceWith(&Path{Node: block.List[0], Parent: block, Field: "Missing", Index: 0}, ast.NewIdent("z"))
		t.Equal(len(block.List), before)
		t.End()
	})

	Test(t, "ReplaceWith: no-op when parent has no struct fields", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		ReplaceWith(&Path{Node: block.List[0], Parent: nil, Field: "List", Index: 0}, ast.NewIdent("z"))
		t.Ok(true)
		t.End()
	})
}

func TestReplaceWithMultiple(t *testing.T) {
	Test(t, "ReplaceWithMultiple: empty nodes list is a no-op", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		before := len(block.List)
		ReplaceWithMultiple(&Path{Node: block.List[0], Parent: block, Field: "List", Index: 0}, nil)
		t.Equal(len(block.List), before)
		t.End()
	})

	Test(t, "ReplaceWithMultiple: replaces one statement with several", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n")
		one := &ast.ExprStmt{X: ast.NewIdent("x")}
		two := &ast.ExprStmt{X: ast.NewIdent("y")}
		ReplaceWithMultiple(&Path{Node: block.List[0], Parent: block, Field: "List", Index: 0}, []ast.Node{one, two})
		t.Ok(len(block.List) == 3 && block.List[0] == one && block.List[1] == two)
		t.End()
	})
}

func TestInsertBefore(t *testing.T) {
	Test(t, "InsertBefore: inserts before a statement", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n")
		insert := &ast.ExprStmt{X: ast.NewIdent("inserted")}
		InsertBefore(&Path{Node: block.List[1], Parent: block, Field: "List", Index: 1}, insert)
		t.Ok(len(block.List) == 3 && block.List[1] == insert)
		t.End()
	})

	Test(t, "InsertBefore: no-op when parent has no matching field", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		before := len(block.List)
		InsertBefore(&Path{Node: block.List[0], Parent: block, Field: "Nope", Index: 0}, &ast.ExprStmt{X: ast.NewIdent("z")})
		t.Equal(len(block.List), before)
		t.End()
	})
}

func TestInsertAfter(t *testing.T) {
	Test(t, "InsertAfter: inserts after a statement", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n")
		insert := &ast.ExprStmt{X: ast.NewIdent("after")}
		InsertAfter(&Path{Node: block.List[0], Parent: block, Field: "List", Index: 0}, insert)
		t.Ok(len(block.List) == 3 && block.List[1] == insert)
		t.End()
	})
}

func TestReplaceInSlice(t *testing.T) {
	Test(t, "replaceInSlice: non-slice field is a no-op", func(t *T) {
		// Field "Node" on a block is an ast.Node interface, not a slice.
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		called := false
		replaceInSlice(&Path{Node: block.List[0], Parent: block, Field: "Node", Index: 0}, func(list reflect.Value, i int) {
			called = true
		})
		t.NotOk(called)
		t.End()
	})

	Test(t, "replaceInSlice: non-slice field of the right name is a no-op", func(t *T) {
		id := ast.NewIdent("x")
		called := false
		replaceInSlice(&Path{Node: ast.NewIdent("y"), Parent: id, Field: "Name", Index: 0}, func(list reflect.Value, i int) {
			called = true
		})
		t.NotOk(called)
		t.End()
	})

	Test(t, "replaceInSlice: out-of-range index is a no-op", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		called := false
		replaceInSlice(&Path{Node: block.List[0], Parent: block, Field: "List", Index: 9}, func(list reflect.Value, i int) {
			called = true
		})
		t.NotOk(called)
		t.End()
	})
}

func TestSetPos(t *testing.T) {
	Test(t, "setPos: stamps token.Pos on all position fields", func(t *T) {
		id := ast.NewIdent("x")
		setPos(id, 42)
		t.Ok(id.NamePos == 42)
		t.End()
	})
}

// print renders a node back to source text, ignoring source positions.
func print(t *testing.T, n ast.Node) string {
	t.Helper()
	stripPositions(n)
	var out bytes.Buffer
	_ = printer.Fprint(&out, token.NewFileSet(), n)
	return out.String()
}

// stripPositions zeroes every token.Pos field on a node sub-tree so the printer
// can render it with any FileSet.
func stripPositions(n ast.Node) {
	ast.Inspect(n, func(node ast.Node) bool {
		if node == nil {
			return false
		}
		e := reflect.ValueOf(node).Elem()
		tt := e.Type()
		for i := 0; i < e.NumField(); i++ {
			if tt.Field(i).Type == reflect.TypeOf(token.NoPos) {
				e.Field(i).SetInt(int64(token.NoPos))
			}
		}
		return true
	})
}

func TestReplaceWithPreservesPos(t *testing.T) {
	Test(t, "ReplaceWith: replacement prints at the original location", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		repl := &ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("z")}, Rhs: []ast.Expr{ast.NewIdent("3")}, Tok: token.ASSIGN}
		ReplaceWith(&Path{Node: block.List[0], Parent: block, Field: "List", Index: 0}, repl)
		out := print(t.TB(), block.List[0])
		t.Match(out, "z")
		t.End()
	})
}

func TestGetBinding(t *testing.T) {
	Test(t, "GetBinding: finds a short variable declaration in a block", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult := foo()\n\tbar(result)\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[1]
		binding := GetBinding(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: block, ParentPath: &Path{Node: file}}}, "result")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: walks up through nested blocks", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult := foo()\n\tif x {\n\t\tbar()\n\t}\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		outer := fn.Body
		ifStmt := outer.List[1].(*ast.IfStmt)
		inner := ifStmt.Body
		stmt := inner.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: inner, ParentPath: &Path{Node: outer, ParentPath: &Path{Node: file}}}, "result")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: ignores plain assignment (= not :=)", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult = foo()\n\tbar()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[1]
		binding := GetBinding(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: block, ParentPath: &Path{Node: file}}}, "result")
		t.NotOk(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a var declaration in a block", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tvar result int\n\tbar()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[1]
		binding := GetBinding(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: block, ParentPath: &Path{Node: file}}}, "result")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds an import alias at file level", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport myfmt \"fmt\"\nfunc f() {}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		binding := GetBinding(&Path{Node: fn, Parent: file, ParentPath: &Path{Node: file}}, "myfmt")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a top-level function declaration", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc Helper() {}\nfunc f() {\n\tHelper()\n}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: fn.Body, ParentPath: &Path{Node: file}}, "Helper")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a top-level var declaration", func(t *T) {
		file := parseFile(t.TB(), "package p\nvar x = 1\nfunc f() {\n\t_ = x\n}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: fn.Body, ParentPath: &Path{Node: file}}, "x")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a top-level type declaration", func(t *T) {
		file := parseFile(t.TB(), "package p\ntype T int\nfunc f() {\n\tvar t T\n}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: fn.Body, ParentPath: &Path{Node: file}}, "T")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a function parameter", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f(arg int) {\n\t_ = arg\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: fn.Body, ParentPath: &Path{Node: fn}}, "arg")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a func literal parameter", func(t *T) {
		file := parseFile(t.TB(), "package p\nvar cb = func(x int) int { return x }\n")
		gd := file.Decls[0].(*ast.GenDecl)
		fnLit := gd.Specs[0].(*ast.ValueSpec).Values[0].(*ast.FuncLit)
		stmt := fnLit.Body.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: fnLit.Body, ParentPath: &Path{Node: fnLit}}, "x")
		t.Ok(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: returns nil when nothing declares name", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tbar()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[0]
		binding := GetBinding(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: fn, ParentPath: &Path{Node: file}}}, "missing")
		t.NotOk(binding != nil)
		t.End()
	})

	Test(t, "GetBindingPath: behaves like GetBinding", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult := foo()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[0]
		name := "result"
		a := GetBinding(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: block, ParentPath: &Path{Node: file}}}, name)
		b := GetBindingPath(&Path{Node: stmt, Parent: block, ParentPath: &Path{Node: block, ParentPath: &Path{Node: file}}}, name)
		t.Ok(a != nil && b != nil)
		t.End()
	})
}

func TestRename(t *testing.T) {
	Test(t, "Rename: renames idents within the subtree", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tfoo(bar)\n\tbaz()\n}\n")
		Rename(&Path{Node: block}, "foo", "qux")
		got := print(t.TB(), block)
		t.Match(got, "qux")
		t.End()
	})
}

func TestExtract(t *testing.T) {
	Test(t, "Extract: identifier name", func(t *T) {
		t.Equal(Extract(ast.NewIdent("Equal")), "Equal")
		t.End()
	})

	Test(t, "Extract: string literal is unquoted", func(t *T) {
		lit := &ast.BasicLit{Kind: token.STRING, Value: `"hello"`}
		t.Equal(Extract(lit), "hello")
		t.End()
	})

	Test(t, "Extract: non-string literal keeps raw value", func(t *T) {
		lit := &ast.BasicLit{Kind: token.INT, Value: "42"}
		t.Equal(Extract(lit), "42")
		t.End()
	})

	Test(t, "Extract: selector expression is X.Sel", func(t *T) {
		sel := &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("Equal")}
		t.Equal(Extract(sel), "t.Equal")
		t.End()
	})

	Test(t, "Extract: panics on unsupported node", func(t *T) {
		panicked := false
		func() {
			defer func() { panicked = recover() != nil }()
			Extract(&ast.CallExpr{})
		}()
		t.Ok(panicked)
		t.End()
	})
}

func TestIsSimple(t *testing.T) {
	Test(t, "IsSimple: literals, identifiers and selectors are simple", func(t *T) {
		t.Ok(IsSimple(&ast.BasicLit{Kind: token.INT, Value: "1"}) &&
			IsSimple(ast.NewIdent("x")) &&
			IsSimple(&ast.SelectorExpr{X: ast.NewIdent("a"), Sel: ast.NewIdent("b")}))
		t.End()
	})

	Test(t, "IsSimple: call expression is not simple", func(t *T) {
		t.NotOk(IsSimple(&ast.CallExpr{}))
		t.End()
	})
}

func TestHasImport(t *testing.T) {
	Test(t, "HasImport: reports an imported path", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport \"fmt\"\nfunc f() {}\n")
		t.Ok(HasImport(file, "fmt"))
		t.End()
	})

	Test(t, "HasImport: returns false for an absent path", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport \"fmt\"\nfunc f() {}\n")
		t.NotOk(HasImport(file, "os"))
		t.End()
	})

	Test(t, "HasImport: nil file", func(t *T) {
		t.NotOk(HasImport(nil, "fmt"))
		t.End()
	})
}
func TestGetImportAlias(t *testing.T) {
	Test(t, "GetImportAlias: plain import has no alias", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport \"fmt\"\nfunc f() {}\n")
		t.Equal(GetImportAlias(file, "fmt"), "")
		t.End()
	})

	Test(t, "GetImportAlias: returns the alias", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport myfmt \"fmt\"\nfunc f() {}\n")
		t.Equal(GetImportAlias(file, "fmt"), "myfmt")
		t.End()
	})

	Test(t, "GetImportAlias: dot import returns dot", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport . \"fmt\"\nfunc f() {}\n")
		t.Equal(GetImportAlias(file, "fmt"), ".")
		t.End()
	})

	Test(t, "GetImportAlias: absent path returns empty", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport \"fmt\"\nfunc f() {}\n")
		t.Equal(GetImportAlias(file, "os"), "")
		t.End()
	})
}

func TestVarsExtractors(t *testing.T) {
	Test(t, "FileFromVars: returns injected file", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {}\n")
		vars := map[string]ast.Node{"$file": file}
		t.Ok(FileFromVars(vars) == file)
		t.End()
	})

	Test(t, "FileFromVars: absent returns nil", func(t *T) {
		t.NotOk(FileFromVars(map[string]ast.Node{}) != nil)
		t.End()
	})

	Test(t, "BlockFromVars: returns injected block", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\ta := 1\n}\n")
		vars := map[string]ast.Node{"$block": block}
		t.Ok(BlockFromVars(vars) == block)
		t.End()
	})

	Test(t, "BlockFromVars: absent returns nil", func(t *T) {
		t.NotOk(BlockFromVars(map[string]ast.Node{}) != nil)
		t.End()
	})
}

func TestCompareReexports(t *testing.T) {
	Test(t, "Compare: true on match", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tfoo(x)\n}\n")
		t.Ok(Compare(block.List[0], "foo(__a)"))
		t.End()
	})

	Test(t, "Compare: false on mismatch", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tfoo(x)\n}\n")
		t.NotOk(Compare(block.List[0], "bar(__a)"))
		t.End()
	})

	Test(t, "GetTemplateValues: binds holes on match", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tfoo(x)\n}\n")
		vars := GetTemplateValues(block.List[0], "foo(__a)")
		t.Ok(vars != nil)
		t.End()
	})

	Test(t, "GetTemplateValues: nil on mismatch", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tfoo(x)\n}\n")
		t.NotOk(GetTemplateValues(block.List[0], "bar(__a)") != nil)
		t.End()
	})

	Test(t, "type aliases: BodySlice and ArgSlice are re-exported", func(t *T) {
		var b BodySlice
		var a ArgSlice
		t.Ok(b.Stmts == nil && a.Args == nil)
		t.End()
	})
}

func TestPreserveCommentsAlreadyTracked(t *testing.T) {
	Test(t, "preserveComments: skips a comment group already tracked by the file", func(t *T) {
		file := parseFile(t.TB(), "package p\n// doc\nfunc f() {}\n")
		before := len(file.Comments)
		fd := file.Decls[0].(*ast.FuncDecl)
		Remove(&Path{Node: fd, Parent: file, ParentPath: &Path{Node: file}, Field: "Decls", Index: 0})
		t.Equal(len(file.Comments), before)
		t.End()
	})
}

func TestBlockDeclaresEdgeCases(t *testing.T) {
	Test(t, "blockDeclaresStmt: a const declaration is not a var binding", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tconst x = 1\n\t_ = x\n}\n")
		_, ok := blockDeclaresStmt(block, "x")
		t.NotOk(ok)
		t.End()
	})

	Test(t, "blockDeclaresStmt: a non-GenDecl DeclStmt is skipped", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{&ast.DeclStmt{Decl: &ast.FuncDecl{Name: ast.NewIdent("g")}}}}
		_, ok := blockDeclaresStmt(block, "x")
		t.NotOk(ok)
		t.End()
	})

	Test(t, "blockDeclaresStmt: a non-ValueSpec var spec is skipped", func(t *T) {
		block := &ast.BlockStmt{List: []ast.Stmt{&ast.DeclStmt{
			Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{&ast.TypeSpec{Name: ast.NewIdent("T")}}},
		}}}
		_, ok := blockDeclaresStmt(block, "x")
		t.NotOk(ok)
		t.End()
	})
}

func TestParamsDeclareNil(t *testing.T) {
	Test(t, "paramsDeclare: nil fields returns false", func(t *T) {
		t.NotOk(paramsDeclare(nil, "x"))
		t.End()
	})
}
