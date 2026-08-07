package apply_fixture_name_to_message

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/types"

	. "github.com/coderaiser/go-tape"
)

// baseSrc has a var Test = CreateTest("remove-skip", nil) declaration and a
// Test message that lacks the "remove-skip: " prefix.
const baseSrc = `package p

import . "coderaiser/indra/internal/test"

var Test = CreateTest("remove-skip", nil)

func f(t *testing.T) {
	Test(t, "report Test.Skip call", func(t *T) {
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

func TestExtractRuleName(t *testing.T) {
	base := parseFile(t, baseSrc)

	Test(t, "extractRuleName: returns rule name from CreateTest", func(t *T) {
		result := extractRuleName(types.Path{Node: base})
		t.Equal(result, "remove-skip")

		t.End()
	})

	Test(t, "extractRuleName: empty when no var Test", func(t *T) {
		src := "package p\nvar Other = CreateTest(\"x\", nil)\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: skip when value is not a call", func(t *T) {
		src := "package p\nvar Test = 5\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: skip when fun is not an ident", func(t *T) {
		src := "package p\nvar Test = x.CreateTest(\"y\")\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: skip when fun is not CreateTest", func(t *T) {
		src := "package p\nvar Test = Foo(\"y\")\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: skip when no args", func(t *T) {
		src := "package p\nvar Test = CreateTest()\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: skip when arg is not a literal", func(t *T) {
		src := "package p\nvar Test = CreateTest(x)\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: empty literal yields empty rule", func(t *T) {
		src := "package p\nvar Test = CreateTest(\"\")\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	Test(t, "extractRuleName: skips index beyond values", func(t *T) {
		src := "package p\nvar X, Test = 5\n"
		result := extractRuleName(types.Path{Node: parseFile(t.TB(), src)})
		t.Equal(result, "")

		t.End()
	})

	// Construct a malformed 1-char BasicLit to cover the len(s) < 2 branch.
	Test(t, "extractRuleName: skips literal shorter than quotes", func(t *T) {
		file := parseFile(t.TB(), "package p\nvar Test = CreateTest(\"\")\n")
		vs := file.Decls[0].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
		call := vs.Values[0].(*ast.CallExpr)
		call.Args[0] = &ast.BasicLit{Kind: token.STRING, Value: "x"}
		result := extractRuleName(types.Path{Node: file})
		t.Equal(result, "")

		t.End()
	})
}

func TestHasMissingPrefix(t *testing.T) {
	Test(t, "hasMissingPrefix: true when message lacks prefix", func(t *T) {
		p := types.Path{Node: parseFile(t.TB(), baseSrc)}
		t.Ok(hasMissingPrefix(p, "remove-skip"))
		t.End()
	})

	Test(t, "hasMissingPrefix: false when message already prefixed", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: report", func(t *T) { t.End() })
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.NotOk(hasMissingPrefix(p, "remove-skip"))

		t.End()
	})

	Test(t, "hasMissingPrefix: later call short-circuits", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "first", func(t *T) { t.End() })
	Test(t, "second", func(t *T) { t.End() })
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.Ok(hasMissingPrefix(p, "remove-skip"))
		t.End()
	})

	Test(t, "hasMissingPrefix: ignores non-matching call shapes", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	t.End()
	Foo(t, "x")
	Test(t)
	Test(t, x)
	Test(t, "remove-skip: ok")
}
`
		p := types.Path{Node: parseFile(t.TB(), src)}
		t.NotOk(hasMissingPrefix(p, "remove-skip"))

		t.End()
	})

	Test(t, "hasMissingPrefix: literal shorter than quotes is ignored", func(t *T) {
		file := parseFile(t.TB(), baseSrc)
		vs := file.Decls[1].(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
		_ = vs
		fn := file.Decls[2].(*ast.FuncDecl)
		call := fn.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		call.Args[1] = &ast.BasicLit{Kind: token.STRING, Value: "x"}
		t.NotOk(hasMissingPrefix(types.Path{Node: file}, "remove-skip"))

		t.End()
	})
}

func TestApplyPrefix(t *testing.T) {
	Test(t, "applyPrefix: prepends missing prefix", func(t *T) {
		file := parseFile(t.TB(), baseSrc)
		applyPrefix(types.Path{Node: file}, "remove-skip")
		fn := file.Decls[2].(*ast.FuncDecl)
		call := fn.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		lit := call.Args[1].(*ast.BasicLit)
		t.Equal(lit.Value, `"remove-skip: report Test.Skip call"`)
		t.End()
	})

	Test(t, "applyPrefix: leaves already-prefixed message unchanged", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	Test(t, "remove-skip: report", func(t *T) { t.End() })
}
`
		file := parseFile(t.TB(), src)
		applyPrefix(types.Path{Node: file}, "remove-skip")
		t.Pass("no panic")
		t.End()
	})

	Test(t, "applyPrefix: ignores non-matching call shapes", func(t *T) {
		src := `package p
var Test = CreateTest("remove-skip", nil)
func f(t *testing.T) {
	t.End()
	Foo(t, "x")
	Test(t)
	Test(t, x)
	Test(t, "remove-skip: ok")
}
`
		file := parseFile(t.TB(), src)
		applyPrefix(types.Path{Node: file}, "remove-skip")
		t.Pass("no panic")
		t.End()
	})
}

func TestFix(t *testing.T) {
	Test(t, "Fix: no-op when no rule name", func(t *T) {
		src := "package p\nfunc f(t *testing.T) { Test(t, \"x\", func(t *T) {}) }\n"
		file := parseFile(t.TB(), src)
		Fix(types.Path{Node: file}, nil)
		t.Pass("no panic")
		t.End()
	})
}
