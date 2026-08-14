package types_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/ast/astutil"

	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

func ident(name string) *ast.Ident { return ast.NewIdent(name) }

func TestFind(t *testing.T) {
	a, b, c := ident("a"), ident("b"), ident("c")
	path := types.Path{Node: c, Stack: []ast.Node{a, b}}

	Test(t, "Path.Find: matches self", func(t *T) {
		_, ok := path.Find(func(p types.Path) bool { return p.Node == c })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.Find: matches ancestor", func(t *T) {
		_, ok := path.Find(func(p types.Path) bool { return p.Node == a })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.Find: returns false when nothing matches", func(t *T) {
		_, ok := path.Find(func(p types.Path) bool { return false })
		t.NotOk(ok)

		t.End()
	})
}

func TestFindParent(t *testing.T) {
	a, b, c := ident("a"), ident("b"), ident("c")
	path := types.Path{Node: c, Stack: []ast.Node{a, b}}

	Test(t, "Path.FindParent: finds immediate parent", func(t *T) {
		_, ok := path.FindParent(func(p types.Path) bool { return p.Node == b })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.FindParent: finds grandparent", func(t *T) {
		_, ok := path.FindParent(func(p types.Path) bool { return p.Node == a })
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.FindParent: does not match self", func(t *T) {
		_, ok := path.FindParent(func(p types.Path) bool { return p.Node == c })
		t.NotOk(ok)

		t.End()
	})

	Test(t, "Path.FindParent: returns false on empty stack", func(t *T) {
		root := types.Path{Node: a, Stack: nil}
		_, ok := root.FindParent(func(p types.Path) bool { return true })
		t.NotOk(ok)

		t.End()
	})

	Test(t, "Path.FindParent: ancestor stack is trimmed correctly", func(t *T) {
		found, _ := path.FindParent(func(p types.Path) bool { return p.Node == b })
		result := len(found.Stack)
		t.Equal(result, 1)

		t.End()
	})
}

func TestParentPath(t *testing.T) {
	a, b := ident("a"), ident("b")
	path := types.Path{Node: b, Stack: []ast.Node{a}}

	Test(t, "Path.ParentPath: returns parent", func(t *T) {
		_, ok := path.ParentPath()
		t.Ok(ok)

		t.End()
	})

	Test(t, "Path.ParentPath: returns false at root", func(t *T) {
		_, ok := types.Path{Node: a, Stack: nil}.ParentPath()
		t.NotOk(ok)

		t.End()
	})

	Test(t, "Path.ParentPath: parent has empty stack at root", func(t *T) {
		parent, _ := path.ParentPath()
		result := len(parent.Stack)
		t.Equal(result, 0)

		t.End()
	})
}

func TestPathReplace(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n\ty := 2\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		if assign, ok := c.Node().(*ast.AssignStmt); ok {
			if len(assign.Lhs) > 0 {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == "x" {
					path := types.Path{Node: assign, Cursor: c}
					path.Replace(&ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("z")},
						Rhs:    []ast.Expr{ast.NewIdent("3")},
						TokPos: assign.TokPos,
						Tok:    token.DEFINE,
					})
					return false
				}
			}
		}
		return true
	}, nil)

	funcDecl := file.Decls[0].(*ast.FuncDecl)

	Test(t, "Path.Replace: keeps list length", func(t *T) {
		result := len(funcDecl.Body.List)
		t.Equal(result, 2)
		t.End()
	})

	Test(t, "Path.Replace: updates first statement lhs", func(t *T) {
		first := funcDecl.Body.List[0].(*ast.AssignStmt)
		result := first.Lhs[0].(*ast.Ident).Name
		t.Equal(result, "z")
		t.End()
	})
}

func TestPathDelete(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n\ty := 2\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		if assign, ok := c.Node().(*ast.AssignStmt); ok {
			if len(assign.Lhs) > 0 {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == "x" {
					path := types.Path{Node: assign, Cursor: c}
					path.Delete()
					return false
				}
			}
		}
		return true
	}, nil)

	funcDecl := file.Decls[0].(*ast.FuncDecl)

	Test(t, "Path.Delete: removes node via cursor", func(t *T) {
		result := len(funcDecl.Body.List)
		t.Equal(result, 1)

		t.End()
	})
}

func TestPathInsertBefore(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		if assign, ok := c.Node().(*ast.AssignStmt); ok {
			if len(assign.Lhs) > 0 {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == "x" {
					path := types.Path{Node: assign, Cursor: c}
					path.InsertBefore(&ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("z")},
						Rhs:    []ast.Expr{ast.NewIdent("3")},
						TokPos: assign.TokPos,
						Tok:    token.DEFINE,
					})
					return false
				}
			}
		}
		return true
	}, nil)

	funcDecl := file.Decls[0].(*ast.FuncDecl)

	Test(t, "Path.InsertBefore: keeps list length", func(t *T) {
		result := len(funcDecl.Body.List)
		t.Equal(result, 2)
		t.End()
	})

	Test(t, "Path.InsertBefore: prepends new statement", func(t *T) {
		first := funcDecl.Body.List[0].(*ast.AssignStmt)
		result := first.Lhs[0].(*ast.Ident).Name
		t.Equal(result, "z")
		t.End()
	})
}

func TestPathInsertAfter(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)

	astutil.Apply(file, func(c *astutil.Cursor) bool {
		if assign, ok := c.Node().(*ast.AssignStmt); ok {
			if len(assign.Lhs) > 0 {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == "x" {
					path := types.Path{Node: assign, Cursor: c}
					path.InsertAfter(&ast.AssignStmt{
						Lhs:    []ast.Expr{ast.NewIdent("z")},
						Rhs:    []ast.Expr{ast.NewIdent("3")},
						TokPos: assign.TokPos,
						Tok:    token.DEFINE,
					})
					return false
				}
			}
		}
		return true
	}, nil)

	funcDecl := file.Decls[0].(*ast.FuncDecl)

	Test(t, "Path.InsertAfter: keeps list length", func(t *T) {
		result := len(funcDecl.Body.List)
		t.Equal(result, 2)
		t.End()
	})

	Test(t, "Path.InsertAfter: appends new statement", func(t *T) {
		second := funcDecl.Body.List[1].(*ast.AssignStmt)
		result := second.Lhs[0].(*ast.Ident).Name
		t.Equal(result, "z")
		t.End()
	})
}

func TestPathCursorMethodsNoOpWhenNil(t *testing.T) {
	path := types.Path{Node: ast.NewIdent("x"), Cursor: nil}

	// Should not panic
	path.Replace(ast.NewIdent("y"))
	path.Delete()
	path.InsertBefore(ast.NewIdent("z"))
	path.InsertAfter(ast.NewIdent("w"))

	Test(t, "Path cursor methods: no-ops when Cursor is nil", func(t *T) {
		t.Pass("no panic")
		t.End()
	})
}

func TestPathStopNoOpOnEnginePath(t *testing.T) {
	path := types.Path{Node: ast.NewIdent("x")}

	// Should not panic — state is nil on engine-constructed paths
	path.Stop()

	Test(t, "Path.Stop: no-op on engine-constructed path", func(t *T) {
		t.Pass("no panic")
		t.End()
	})
}

func TestPathSkipNoOpOnEnginePath(t *testing.T) {
	path := types.Path{Node: ast.NewIdent("x")}

	// Should not panic — state is nil on engine-constructed paths
	path.Skip()

	Test(t, "Path.Skip: no-op on engine-constructed path", func(t *T) {
		t.Pass("no panic")
		t.End()
	})
}

func TestPathTraversePatternKey(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n\tt.Ok(x)\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	funcDecl := file.Decls[0].(*ast.FuncDecl)
	p := types.Path{Node: funcDecl.Body}

	Test(t, "Path.Traverse: pattern key fires on matching ExprStmt", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"t.Equal(__a, __b)": func(_ types.Path) { count++ },
		})
		t.Equal(count, 1)
		t.End()
	})

	Test(t, "Path.Traverse: pattern key does not fire on non-matching node", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"t.Equal(__a, __b, __c)": func(_ types.Path) { count++ },
		})
		t.Equal(count, 0)
		t.End()
	})

	Test(t, "Path.Traverse: type key fires on all ExprStmt nodes", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"*ast.ExprStmt": func(_ types.Path) { count++ },
		})
		t.Equal(count, 2)
		t.End()
	})

	Test(t, "Path.Traverse: pattern key fires once when one stmt matches", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"t.Equal(__a, __b)": func(_ types.Path) { count++ },
		})
		t.Equal(count, 1)
		t.End()
	})

	Test(t, "Path.Traverse: pattern key Stop halts after first match", func(t *T) {
		src2 := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n\tt.Equal(c, d)\n}\n"
		fset2 := token.NewFileSet()
		file2, _ := parser.ParseFile(fset2, "", src2, 0)
		fn2 := file2.Decls[0].(*ast.FuncDecl)
		p2 := types.Path{Node: fn2.Body}
		count := 0
		p2.Traverse(map[string]func(types.Path){
			"t.Equal(__a, __b)": func(child types.Path) {
				count++
				child.Stop()
			},
		})
		t.Equal(count, 1)
		t.End()
	})

	Test(t, "Path.Traverse: pattern key Skip skips matched node children", func(t *T) {
		src3 := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
		fset3 := token.NewFileSet()
		file3, _ := parser.ParseFile(fset3, "", src3, 0)
		fn3 := file3.Decls[0].(*ast.FuncDecl)
		p3 := types.Path{Node: fn3.Body}
		callExprs := 0
		p3.Traverse(map[string]func(types.Path){
			"t.Equal(__a, __b)": func(child types.Path) { child.Skip() },
			"*ast.CallExpr":     func(_ types.Path) { callExprs++ },
		})
		// the t.Equal(a, b) CallExpr is a child of the skipped ExprStmt
		t.Equal(callExprs, 0)
		t.End()
	})
}

func TestPathTraverse(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n\treturn\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	funcDecl := file.Decls[0].(*ast.FuncDecl)
	p := types.Path{Node: funcDecl.Body}

	Test(t, "Path.Traverse: visits matching node type", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"*ast.ReturnStmt": func(_ types.Path) { count++ },
		})
		t.Equal(count, 1)
		t.End()
	})

	Test(t, "Path.Traverse: skips non-matching node types", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"*ast.ReturnStmt": func(_ types.Path) { count++ },
		})
		// only ReturnStmt, not AssignStmt
		t.Equal(count, 1)
		t.End()
	})

	Test(t, "Path.Traverse: visited path carries parent in stack", func(t *T) {
		var stack []ast.Node
		p.Traverse(map[string]func(types.Path){
			"*ast.ReturnStmt": func(child types.Path) {
				stack = child.Stack
			},
		})
		t.Ok(len(stack) > 0)
		t.End()
	})

	Test(t, "Path.Traverse: no-op for empty visitors map", func(t *T) {
		// must not panic
		p.Traverse(map[string]func(types.Path){})
		t.Ok(true)
		t.End()
	})
}

func TestPathStop(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n\ty := 2\n\tz := 3\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	funcDecl := file.Decls[0].(*ast.FuncDecl)
	p := types.Path{Node: funcDecl.Body}

	Test(t, "Path.Stop: halts traversal after first match", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"*ast.AssignStmt": func(child types.Path) {
				count++
				child.Stop()
			},
		})
		t.Equal(count, 1)
		t.End()
	})

	Test(t, "Path.Stop: does not affect a later independent Traverse call", func(t *T) {
		count := 0
		p.Traverse(map[string]func(types.Path){
			"*ast.AssignStmt": func(child types.Path) { child.Stop() },
		})
		p.Traverse(map[string]func(types.Path){
			"*ast.AssignStmt": func(child types.Path) { count++ },
		})
		// the second call visits all three assignments independently
		t.Equal(count, 3)
		t.End()
	})
}

func TestLintContract(t *testing.T) {
	Test(t, "types: Lint returns no error", func(t *T) {
		lint := types.Lint(func(src []byte, fix bool, plugins []any) (types.LintResult, error) {
			return types.LintResult{Out: src, Places: nil}, nil
		})
		_, err := lint([]byte("x"), false, nil)
		t.NotOk(err)

		t.End()
	})

	Test(t, "types: LintResult carries Out", func(t *T) {
		lint := types.Lint(func(src []byte, fix bool, plugins []any) (types.LintResult, error) {
			return types.LintResult{Out: src, Places: nil}, nil
		})
		result, _ := lint([]byte("x"), false, nil)
		got := string(result.Out)
		t.Equal(got, "x")
		t.End()
	})
}

func TestPathSkip(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tif true {\n\t\ty := 2\n\t}\n\tz := 3\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	funcDecl := file.Decls[0].(*ast.FuncDecl)
	p := types.Path{Node: funcDecl.Body}

	Test(t, "Path.Skip: skips children of current node only", func(t *T) {
		assigns := 0
		p.Traverse(map[string]func(types.Path){
			"*ast.IfStmt":     func(child types.Path) { child.Skip() },
			"*ast.AssignStmt": func(_ types.Path) { assigns++ },
		})
		// y := 2 is inside the skipped if; z := 3 is a sibling and is visited
		t.Equal(assigns, 1)
		t.End()
	})
}
