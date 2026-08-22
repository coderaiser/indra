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
		t.NotOk(result)

		t.End()
	})
}

func TestPrevSibling(t *testing.T) {
	src := "package p\n\nfunc f() {\n\ta := 1\n\tb := 2\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	fn := file.Decls[0].(*ast.FuncDecl)
	block := fn.Body

	Test(t, "Path.PrevSibling: returns previous statement", func(t *T) {
		p := types.Path{Node: block.List[1], Stack: []ast.Node{file, fn, block}}
		prev, ok := p.PrevSibling()
		t.Ok(ok && prev.Node == block.List[0])
		t.End()
	})

	Test(t, "Path.PrevSibling: first statement has no prev sibling", func(t *T) {
		p := types.Path{Node: block.List[0], Stack: []ast.Node{file, fn, block}}
		_, ok := p.PrevSibling()
		t.NotOk(ok)
		t.End()
	})

	Test(t, "Path.PrevSibling: empty stack returns false", func(t *T) {
		_, ok := types.Path{Node: block.List[0], Stack: nil}.PrevSibling()
		t.NotOk(ok)
		t.End()
	})

	Test(t, "Path.PrevSibling: non-list parent returns false", func(t *T) {
		p := types.Path{Node: ast.NewIdent("x"), Stack: []ast.Node{(ast.Node)(&ast.ExprStmt{X: ast.NewIdent("y")})}}
		_, ok := p.PrevSibling()
		t.NotOk(ok)
		t.End()
	})

	Test(t, "Path.PrevSibling: node not found in parent list returns false", func(t *T) {
		p := types.Path{Node: ast.NewIdent("x"), Stack: []ast.Node{file, fn, block}}
		_, ok := p.PrevSibling()
		t.NotOk(ok)
		t.End()
	})

	Test(t, "Path.PrevSibling: file-level declarations have prev siblings", func(t *T) {
		// file.Decls[0] is the single func; give it a sibling to observe reordering
		importDecl := &ast.GenDecl{Tok: token.IMPORT}
		file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
		p := types.Path{Node: file.Decls[1], Stack: []ast.Node{file}}
		prev, ok := p.PrevSibling()
		t.Ok(ok && prev.Node == importDecl)
		t.End()
	})
}

func TestCompareHelpers(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "", src, 0)
	stmt := file.Decls[0].(*ast.FuncDecl).Body.List[0]

	Test(t, "Compare: matches a statement against a pattern", func(t *T) {
		t.Ok(types.Compare(stmt, "__a.Equal(__b, __c)"))
		t.End()
	})

	Test(t, "Compare: non-matching pattern is rejected", func(t *T) {
		t.NotOk(types.Compare(stmt, "t.End()"))
		t.End()
	})

	Test(t, "GetTemplateValues: binds holes from a matching pattern", func(t *T) {
		vars := types.GetTemplateValues(stmt, "__a.Equal(__b, __c)")
		t.Ok(vars != nil && vars["__b"].(*ast.Ident).Name == "a")
		t.End()
	})

	Test(t, "GetTemplateValues: returns nil for a non-match", func(t *T) {
		t.NotOk(types.GetTemplateValues(stmt, "t.End()") != nil)
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
		t.NotOk(count)

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
		t.NotOk(
			// the t.Equal(a, b) CallExpr is a child of the skipped ExprStmt
			callExprs)

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

func TestBabelReexports(t *testing.T) {
	sel := &ast.SelectorExpr{X: ident("a"), Sel: ident("b")}
	arrayLit := &ast.CompositeLit{Type: &ast.ArrayType{Elt: ident("int")}}
	structLit := &ast.CompositeLit{Type: ident("T")}
	anonLit := &ast.CompositeLit{Elts: []ast.Expr{ident("a")}}
	call := &ast.CallExpr{Fun: ident("f")}
	stmt := &ast.ExprStmt{X: call}
	intLit := &ast.BasicLit{Kind: token.INT, Value: "1"}

	Test(t, "re-export: IsIdent accepts an identifier", func(t *T) {
		t.Ok(types.IsIdent(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsIdent rejects a literal", func(t *T) {
		t.NotOk(types.IsIdent(intLit))
		t.End()
	})
	Test(t, "re-export: IsCallExpr accepts a call", func(t *T) {
		t.Ok(types.IsCallExpr(call))
		t.End()
	})
	Test(t, "re-export: IsCallExpr rejects an ident", func(t *T) {
		t.NotOk(types.IsCallExpr(ident("f")))
		t.End()
	})
	Test(t, "re-export: IsSelector accepts a selector", func(t *T) {
		t.Ok(types.IsSelector(sel))
		t.End()
	})
	Test(t, "re-export: IsSelector rejects an ident", func(t *T) {
		t.NotOk(types.IsSelector(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsCompositeLit accepts a composite literal", func(t *T) {
		t.Ok(types.IsCompositeLit(structLit))
		t.End()
	})
	Test(t, "re-export: IsCompositeLit rejects an ident", func(t *T) {
		t.NotOk(types.IsCompositeLit(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsArrayExpr accepts a slice composite", func(t *T) {
		t.Ok(types.IsArrayExpr(arrayLit))
		t.End()
	})
	Test(t, "re-export: IsArrayExpr rejects anonymous and non-slice composites", func(t *T) {
		t.NotOk(types.IsArrayExpr(anonLit) || types.IsArrayExpr(structLit) || types.IsArrayExpr(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsObjectExpr accepts a named struct composite", func(t *T) {
		t.Ok(types.IsObjectExpr(structLit))
		t.End()
	})
	Test(t, "re-export: IsObjectExpr accepts a qualified struct composite", func(t *T) {
		t.Ok(types.IsObjectExpr(&ast.CompositeLit{Type: sel}))
		t.End()
	})
	Test(t, "re-export: IsObjectExpr rejects slice, anonymous and non-composites", func(t *T) {
		t.NotOk(types.IsObjectExpr(arrayLit) || types.IsObjectExpr(anonLit) || types.IsObjectExpr(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsFuncLit accepts a func literal", func(t *T) {
		t.Ok(types.IsFuncLit(&ast.FuncLit{Type: &ast.FuncType{}}))
		t.End()
	})
	Test(t, "re-export: IsFuncLit rejects an ident", func(t *T) {
		t.NotOk(types.IsFuncLit(ident("f")))
		t.End()
	})
	Test(t, "re-export: IsBasicLit accepts a literal", func(t *T) {
		t.Ok(types.IsBasicLit(intLit))
		t.End()
	})
	Test(t, "re-export: IsBasicLit rejects an ident", func(t *T) {
		t.NotOk(types.IsBasicLit(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsStatement accepts a statement", func(t *T) {
		t.Ok(types.IsStatement(stmt))
		t.End()
	})
	Test(t, "re-export: IsStatement rejects an expression", func(t *T) {
		t.NotOk(types.IsStatement(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsFile accepts a file", func(t *T) {
		t.Ok(types.IsFile(&ast.File{}))
		t.End()
	})
	Test(t, "re-export: IsFile rejects a non-file", func(t *T) {
		t.NotOk(types.IsFile(ident("a")))
		t.End()
	})
	Test(t, "re-export: IsBoolLit matches true and false literals", func(t *T) {
		t.Ok(types.IsBoolLit(ident("true"), true) && types.IsBoolLit(ident("false"), false))
		t.End()
	})
	Test(t, "re-export: IsBoolLit rejects mismatches and non-bool idents", func(t *T) {
		t.NotOk(types.IsBoolLit(ident("true"), false) || types.IsBoolLit(ident("x"), true))
		t.End()
	})
	Test(t, "re-export: IsBoolLit rejects non-identifiers", func(t *T) {
		t.NotOk(types.IsBoolLit(intLit, true))
		t.End()
	})
}

func TestOptionsStringSlice(t *testing.T) {
	Test(t, "Options.StringSlice: missing key returns nil", func(t *T) {
		t.NotOk(types.Options{}.StringSlice("allowed"))

		t.End()
	})

	Test(t, "Options.StringSlice: string becomes a single-element slice", func(t *T) {
		got := types.Options{"allowed": "Suite"}.StringSlice("allowed")
		t.Ok(len(got) == 1 && got[0] == "Suite")
		t.End()
	})

	Test(t, "Options.StringSlice: string slice passes through", func(t *T) {
		got := types.Options{"allowed": []string{"a", "b"}}.StringSlice("allowed")
		t.Ok(len(got) == 2 && got[0] == "a" && got[1] == "b")
		t.End()
	})

	Test(t, "Options.StringSlice: any slice of strings is collected", func(t *T) {
		got := types.Options{"allowed": []any{"a", "b"}}.StringSlice("allowed")
		t.Ok(len(got) == 2 && got[1] == "b")
		t.End()
	})

	Test(t, "Options.StringSlice: non-string any elements are skipped", func(t *T) {
		got := types.Options{"allowed": []any{"a", 1, "c"}}.StringSlice("allowed")
		t.Ok(len(got) == 2 && got[0] == "a" && got[1] == "c")
		t.End()
	})

	Test(t, "Options.StringSlice: unsupported value type returns nil", func(t *T) {
		t.NotOk(types.Options{"allowed": 42}.StringSlice("allowed"))

		t.End()
	})
}
