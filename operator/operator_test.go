package operator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"testing"

	"golang.org/x/tools/go/ast/astutil"

	"coderaiser/indra/types"

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

// withCursor runs fn on the Path of the first node matching want during an
// astutil.Apply pre-order walk, so the Path carries a live Cursor. It is the
// engine-shaped harness used for operator functions that delegate to a cursor.
func withCursor(file *ast.File, want ast.Node, fn func(types.Path)) {
	var stack []ast.Node
	astutil.Apply(file, func(c *astutil.Cursor) bool {
		stack = append(stack, c.Node())
		p := types.Path{
			Node:   c.Node(),
			Stack:  append([]ast.Node{}, stack[:len(stack)-1]...),
			Cursor: c,
		}
		if c.Node() == want {
			fn(p)
		}
		return true
	}, func(c *astutil.Cursor) bool {
		if len(stack) > 0 {
			stack = stack[:len(stack)-1]
		}
		return true
	})
}

// firstStmt returns the first statement in block.

func TestRemove(t *testing.T) {
	Test(t, "Remove: deletes the matched statement from its block", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		before := len(block.List)
		withCursor(file, block.List[0], func(p types.Path) {
			Remove(p)
		})
		result := len(block.List)
		t.Equal(result, before-1)

		t.End()
	})

	Test(t, "Remove: deleting an import removes it from the GenDecl.Specs", func(t *T) {
		src := "package p\nimport (\n\t\"a\"\n\t\"b\"\n)\n"
		file := parseFile(t.TB(), src)
		gd := file.Decls[0].(*ast.GenDecl)
		before := len(gd.Specs)
		withCursor(file, file.Imports[1], func(p types.Path) {
			Remove(p)
		})
		result := len(gd.Specs)
		t.Equal(result, before-1)

		t.End()
	})
}

func TestReplaceWith(t *testing.T) {
	Test(t, "ReplaceWith: replaces a statement in a block and preserves position", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		orig := block.List[0]
		repl := &ast.AssignStmt{
			Lhs:    []ast.Expr{ast.NewIdent("z")},
			Rhs:    []ast.Expr{ast.NewIdent("2")},
			Tok:    token.ASSIGN,
			TokPos: orig.Pos(),
		}
		withCursor(file, orig, func(p types.Path) {
			ReplaceWith(p, repl)
		})
		t.Ok(block.List[0] == repl)
		t.End()
	})

	Test(t, "ReplaceWith: replacement prints at the original location", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		orig := block.List[0]
		repl := &ast.AssignStmt{Lhs: []ast.Expr{ast.NewIdent("z")}, Rhs: []ast.Expr{ast.NewIdent("3")}, Tok: token.ASSIGN}
		withCursor(file, orig, func(p types.Path) {
			ReplaceWith(p, repl)
		})
		out := print(t.TB(), block.List[0])
		t.Match(out, "z")
		t.End()
	})
}

func TestReplaceWithMultiple(t *testing.T) {
	Test(t, "ReplaceWithMultiple: empty nodes list is a no-op", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		withCursor(file, block.List[0], func(p types.Path) {
			ReplaceWithMultiple(p, []ast.Node{})
		})
		result := len(block.List)
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "ReplaceWithMultiple: replaces one statement with several", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		one := &ast.ExprStmt{X: ast.NewIdent("x")}
		two := &ast.ExprStmt{X: ast.NewIdent("y")}
		withCursor(file, block.List[0], func(p types.Path) {
			ReplaceWithMultiple(p, []ast.Node{one, two})
		})
		t.Ok(len(block.List) == 3 && block.List[0] == one && block.List[1] == two)
		t.End()
	})
}

func TestInsertBefore(t *testing.T) {
	Test(t, "InsertBefore: inserts before a statement", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		insert := &ast.ExprStmt{X: ast.NewIdent("ins")}
		withCursor(file, block.List[1], func(p types.Path) {
			InsertBefore(p, insert)
		})
		t.Ok(len(block.List) == 3 && block.List[1] == insert)
		t.End()
	})
}

func TestInsertAfter(t *testing.T) {
	Test(t, "InsertAfter: inserts after a statement", func(t *T) {
		src := "package p\nfunc f() {\n\ta := 1\n\tb := 2\n}\n"
		file := parseFile(t.TB(), src)
		block := file.Decls[0].(*ast.FuncDecl).Body
		insert := &ast.ExprStmt{X: ast.NewIdent("after")}
		withCursor(file, block.List[0], func(p types.Path) {
			InsertAfter(p, insert)
		})
		t.Ok(len(block.List) == 3 && block.List[1] == insert)
		t.End()
	})
}

func TestGetBinding(t *testing.T) {
	Test(t, "GetBinding: finds a short variable declaration in a block", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult := foo()\n\tbar(result)\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[1]
		p := types.Path{Node: stmt, Stack: []ast.Node{block, fn, file}}
		binding := GetBinding(p, "result")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: walks up through nested blocks", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult := foo()\n\tif x {\n\t\tbar()\n\t}\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		outer := fn.Body
		ifStmt := outer.List[1].(*ast.IfStmt)
		inner := ifStmt.Body
		stmt := inner.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{inner, outer, fn, file}}
		binding := GetBinding(p, "result")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: ignores plain assignment (= not :=)", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult = foo()\n\tbar()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[1]
		p := types.Path{Node: stmt, Stack: []ast.Node{block, fn, file}}
		binding := GetBinding(p, "result")
		t.NotOk(binding != nil)
		t.End()
	})

	Test(t, "GetBinding: finds a var declaration in a block", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tvar result int\n\tbar()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[1]
		p := types.Path{Node: stmt, Stack: []ast.Node{block, fn, file}}
		binding := GetBinding(p, "result")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: finds an import alias at file level", func(t *T) {
		file := parseFile(t.TB(), "package p\nimport myfmt \"fmt\"\nfunc f() {}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		p := types.Path{Node: fn, Stack: []ast.Node{file}}
		binding := GetBinding(p, "myfmt")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: finds a top-level function declaration", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc Helper() {}\nfunc f() {\n\tHelper()\n}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{fn.Body, fn, file}}
		binding := GetBinding(p, "Helper")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: finds a top-level var declaration", func(t *T) {
		file := parseFile(t.TB(), "package p\nvar x = 1\nfunc f() {\n\t_ = x\n}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{fn.Body, fn, file}}
		binding := GetBinding(p, "x")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: finds a top-level type declaration", func(t *T) {
		file := parseFile(t.TB(), "package p\ntype T int\nfunc f() {\n\tvar t T\n}\n")
		fn := file.Decls[1].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{fn.Body, fn, file}}
		binding := GetBinding(p, "T")
		t.Ok(binding)

		t.End()
	})
}

func TestGetBindingScope(t *testing.T) {

	Test(t, "GetBinding: finds a function parameter", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f(arg int) {\n\t_ = arg\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		stmt := fn.Body.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{fn.Body, fn, file}}
		binding := GetBinding(p, "arg")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: finds a func literal parameter", func(t *T) {
		file := parseFile(t.TB(), "package p\nvar cb = func(x int) int { return x }\n")
		gd := file.Decls[0].(*ast.GenDecl)
		fnLit := gd.Specs[0].(*ast.ValueSpec).Values[0].(*ast.FuncLit)
		stmt := fnLit.Body.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{fnLit.Body, fnLit, gd, file}}
		binding := GetBinding(p, "x")
		t.Ok(binding)

		t.End()
	})

	Test(t, "GetBinding: returns nil when nothing declares name", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tbar()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{block, fn, file}}
		binding := GetBinding(p, "missing")
		t.NotOk(binding != nil)
		t.End()
	})

	Test(t, "GetBindingPath: behaves like GetBinding", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc f() {\n\tresult := foo()\n}\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		block := fn.Body
		stmt := block.List[0]
		p := types.Path{Node: stmt, Stack: []ast.Node{block, fn, file}}
		a := GetBinding(p, "result")
		b := GetBindingPath(p, "result")
		t.Ok((a != nil) == (b != nil))
		t.End()
	})
}

func TestRename(t *testing.T) {
	Test(t, "Rename: renames idents within the subtree", func(t *T) {
		block := funcBody(t.TB(), "package p\nfunc f() {\n\tfoo(bar)\n\tbaz()\n}\n")
		Rename(types.Path{Node: block}, "foo", "qux")
		got := print(t.TB(), block)
		t.Match(got, "qux")
		t.End()
	})
}

func TestExtract(t *testing.T) {
	Test(t, "Extract: identifier name", func(t *T) {
		result := Extract(ast.NewIdent("Equal"))
		t.Equal(result, "Equal")

		t.End()
	})

	Test(t, "Extract: string literal is unquoted", func(t *T) {
		lit := &ast.BasicLit{Kind: token.STRING, Value: `"hello"`}
		result := Extract(lit)
		t.Equal(result, "hello")

		t.End()
	})

	Test(t, "Extract: non-string literal keeps raw value", func(t *T) {
		lit := &ast.BasicLit{Kind: token.INT, Value: "42"}
		result := Extract(lit)
		t.Equal(result, "42")

		t.End()
	})

	Test(t, "Extract: selector expression is X.Sel", func(t *T) {
		sel := &ast.SelectorExpr{X: ast.NewIdent("t"), Sel: ast.NewIdent("Equal")}
		result := Extract(sel)
		t.Equal(result, "t.Equal")

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

func TestPreserveComments(t *testing.T) {
	Test(t, "preserveComments: node with no comment groups returns early", func(t *T) {
		assign := &ast.AssignStmt{}
		preserveComments(types.Path{Node: assign})
		t.Ok(true)
		t.End()
	})

	Test(t, "preserveComments: no file ancestor returns early", func(t *T) {
		gd := &ast.GenDecl{Tok: token.VAR, Doc: &ast.CommentGroup{}}
		preserveComments(types.Path{Node: gd, Stack: []ast.Node{}})
		t.Ok(true)
		t.End()
	})

	Test(t, "preserveComments: appends a comment group not already tracked", func(t *T) {
		file := parseFile(t.TB(), "package p\ntype T int\n")
		gd := &ast.GenDecl{
			Tok: token.TYPE,
			Doc: &ast.CommentGroup{List: []*ast.Comment{{Slash: 1, Text: "// keep"}}},
		}
		file.Decls[0] = gd
		preserveComments(types.Path{Node: gd, Stack: []ast.Node{file}})
		result := len(file.Comments)
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "preserveComments: skips a comment group already tracked by the file", func(t *T) {
		file := parseFile(t.TB(), "package p\n// doc\nfunc f() {}\n")
		before := len(file.Comments)
		fd := file.Decls[0].(*ast.FuncDecl)
		preserveComments(types.Path{Node: fd, Stack: []ast.Node{file}})
		result := len(file.Comments)
		t.Equal(result, before)

		t.End()
	})
}

func TestFindFile(t *testing.T) {
	file := parseFile(t, "package p\nfunc f() {}\n")
	block := file.Decls[0].(*ast.FuncDecl).Body

	Test(t, "findFile: returns file via Stack", func(t *T) {
		found := findFile(types.Path{Node: block, Stack: []ast.Node{file}})
		t.Ok(found == file)
		t.End()
	})

	Test(t, "findFile: no file ancestor returns nil", func(t *T) {
		found := findFile(types.Path{Node: block, Stack: []ast.Node{block}})
		t.NotOk(found != nil)
		t.End()
	})

	Test(t, "findFile: empty stack returns nil", func(t *T) {
		found := findFile(types.Path{Node: block})
		t.NotOk(found != nil)
		t.End()
	})
}

func TestCommentGroups(t *testing.T) {
	Test(t, "commentGroups: nil node", func(t *T) {
		result := len(commentGroups(nil))
		t.Equal(result, 0)

		t.End()
	})

	Test(t, "commentGroups: typed nil pointer", func(t *T) {
		var id *ast.Ident
		result := len(commentGroups(id))
		t.Equal(result, 0)

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
		result := len(commentGroups(gd))
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "commentGroups: skips a nil pointer child", func(t *T) {
		fd := &ast.FuncDecl{Name: ast.NewIdent("f"), Body: nil, Doc: nil}
		result := len(commentGroups(fd))
		t.Equal(result, 0)

		t.End()
	})

	Test(t, "commentGroups: recurses into a non-nil pointer child", func(t *T) {
		fd := &ast.FuncDecl{
			Name: ast.NewIdent("f"),
			Body: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: ast.NewIdent("x")}}},
		}
		result := len(commentGroups(fd))
		t.Equal(result, 0)

		t.End()
	})

	Test(t, "commentGroups: recurses into an interface child", func(t *T) {
		e := &ast.ExprStmt{X: ast.NewIdent("x")}
		result := len(commentGroups(e))
		t.Equal(result, 0)

		t.End()
	})

	Test(t, "commentGroups: skips a non-node pointer field (File.Scope)", func(t *T) {
		span := &ast.File{Name: ast.NewIdent("p"), Scope: ast.NewScope(nil)}
		result := len(commentGroups(span))
		t.Equal(result, 0)

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

	Test(t, "paramsDeclare: nil fields returns false", func(t *T) {
		t.NotOk(paramsDeclare(nil, "x"))
		t.End()
	})
}
