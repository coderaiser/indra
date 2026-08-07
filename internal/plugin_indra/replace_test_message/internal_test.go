package replace_test_message

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

const baseSrc = `package p

var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "remove-skip: report", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
		t.End()
	})
}
`

func parseFile(t *testing.T, src string) *ast.File {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "f.go", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file
}

func TestAfterSeparator(t *testing.T) {
	Test(t, "afterSeparator: returns part after first separator", func(t *T) {
		result := afterSeparator("remove-skip: report remove-skip")
		t.Equal(result, "report remove-skip")

		t.End()
	})

	Test(t, "afterSeparator: returns whole when no separator", func(t *T) {
		result := afterSeparator("just a message")
		t.Equal(result, "just a message")

		t.End()
	})
}

// cb extracts a callback function body as a Path for extractFixtureName tests.
func cb(t *testing.T, src string) types.Path {
	t.Helper()
	file := parseFile(t, src)
	fn := file.Decls[0].(*ast.FuncDecl)
	return types.Path{Node: fn.Body}
}

func TestExtractFixtureName(t *testing.T) {
	Test(t, "extractFixtureName: returns Report arg", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.Report(\"x\", \"y\") }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "x")

		t.End()
	})

	Test(t, "extractFixtureName: returns Transform arg", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.Transform(\"x\") }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "x")

		t.End()
	})

	Test(t, "extractFixtureName: returns NoReport arg", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.NoReport(\"x\") }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "x")

		t.End()
	})

	Test(t, "extractFixtureName: returns NoTransform arg", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.NoTransform(\"x\") }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "x")

		t.End()
	})

	Test(t, "extractFixtureName: stops at second call", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.Report(\"x\", \"y\"); t.End() }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "x")

		t.End()
	})

	Test(t, "extractFixtureName: skips non-selector call", func(t *T) {
		src := "package p\nfunc cb(t *T) { Foo(\"x\") }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractFixtureName: skips selector on non-t", func(t *T) {
		src := "package p\nfunc cb(t *T) { x.Report(\"y\") }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractFixtureName: skips non-fixture method", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.End() }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractFixtureName: skips call with no args", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.Report() }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractFixtureName: skips non-literal arg", func(t *T) {
		src := "package p\nfunc cb(t *T) { t.Report(x) }\n"
		result := extractFixtureName(cb(t.TB(), src))
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractFixtureName: skips malformed short literal", func(t *T) {
		file := parseFile(t.TB(), "package p\nfunc cb(t *T) { t.Report(\"\", \"y\") }\n")
		fn := file.Decls[0].(*ast.FuncDecl)
		call := fn.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		call.Args[0] = &ast.BasicLit{Kind: token.STRING, Value: "x"}
		result := extractFixtureName(types.Path{Node: fn.Body})
		t.Equal(result, "")

		t.End()
	})
}

func TestHasMissingFixtureName(t *testing.T) {
	Test(t, "hasMissingFixtureName: true when fixture name missing after separator", func(t *T) {
		p := types.Path{Node: parseFile(t.TB(), baseSrc)}
		t.Ok(hasMissingFixtureName(p))
		t.End()
	})

	Test(t, "hasMissingFixtureName: false when already correct", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: report remove-skip", func(t *T) {
		t.Report("remove-skip", "remove Test.Skip call")
	})
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.NotOk(hasMissingFixtureName(p))

		t.End()
	})

	Test(t, "hasMissingFixtureName: false when callback has no fixture call", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: some message", func(t *T) { t.End() })
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.NotOk(hasMissingFixtureName(p))

		t.End()
	})

	Test(t, "hasMissingFixtureName: short-circuits on second test", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: report", func(t *T) { t.Report("remove-skip", "x") })
	Test(t, "remove-skip: more", func(t *T) { t.Report("remove-skip", "x") })
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.Ok(hasMissingFixtureName(p))
		t.End()
	})

	Test(t, "hasMissingFixtureName: ignores non-matching shapes", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	t.End()
	Foo(t, "x")
	Test(t, "x")
	Test(t, x, func(t *T) {})
	Test(t, "remove-skip: report remove-skip", 5)
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.NotOk(hasMissingFixtureName(p))

		t.End()
	})

	Test(t, "hasMissingFixtureName: malformed short literal is ignored", func(t *T) {
		file := parseFile(t.TB(), baseSrc)
		fn := file.Decls[1].(*ast.FuncDecl)
		call := fn.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		call.Args[1] = &ast.BasicLit{Kind: token.STRING, Value: "x"}
		t.NotOk(hasMissingFixtureName(types.Path{Node: file}))

		t.End()
	})
}

func TestApplyFixtureNames(t *testing.T) {
	Test(t, "applyFixtureNames: appends fixture name when missing", func(t *T) {
		file := parseFile(t.TB(), baseSrc)
		applyFixtureNames(types.Path{Node: file})
		fn := file.Decls[1].(*ast.FuncDecl)
		call := fn.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		lit := call.Args[1].(*ast.BasicLit)
		t.Equal(lit.Value, `"remove-skip: report remove-skip"`)
		t.End()
	})

	Test(t, "applyFixtureNames: leaves already-correct message unchanged", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: report remove-skip", func(t *T) { t.Report("remove-skip", "x") })
}
`
		file := parseFile(t.TB(), src)
		applyFixtureNames(types.Path{Node: file})
		t.Pass("no panic")
		t.End()
	})

	Test(t, "applyFixtureNames: ignores non-matching shapes", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	t.End()
	Foo(t, "x")
	Test(t, "x")
	Test(t, x, func(t *T) {})
	Test(t, "remove-skip: report remove-skip", 5)
}
`
		file := parseFile(t.TB(), src)
		applyFixtureNames(types.Path{Node: file})
		t.Pass("no panic")
		t.End()
	})

	Test(t, "applyFixtureNames: skips malformed short message literal", func(t *T) {
		file := parseFile(t.TB(), baseSrc)
		fn := file.Decls[1].(*ast.FuncDecl)
		call := fn.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		call.Args[1] = &ast.BasicLit{Kind: token.STRING, Value: "x"}
		applyFixtureNames(types.Path{Node: file})
		t.Pass("no panic")
		t.End()
	})

	Test(t, "applyFixtureNames: skips when callback has no fixture call", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: some", func(t *T) { t.End() })
}
`
		file := parseFile(t.TB(), src)
		applyFixtureNames(types.Path{Node: file})
		t.Pass("no panic")
		t.End()
	})
}
