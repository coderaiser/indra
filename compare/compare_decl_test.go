package compare_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/compare"

	. "github.com/coderaiser/go-tape"
)

func parseDecl(src string) ast.Decl {
	f, err := parser.ParseFile(token.NewFileSet(), "", "package p\n"+src, 0)
	if err != nil || len(f.Decls) == 0 {
		return nil
	}
	return f.Decls[0]
}

func TestCompareDecl(t *testing.T) {
	decl := parseDecl(`func Match() Matcher { return Matcher{"k": nil} }`)

	Test(t, "CompareDecl: matching decl returns non-nil vars", func(t *T) {
		vars := compare.CompareDecl(decl, `func Match() Matcher { return Matcher{__a: nil} }`)
		t.Ok(vars)

		t.End()
	})

	Test(t, "CompareDecl: matching decl binds hole", func(t *T) {
		vars := compare.CompareDecl(decl, `func Match() Matcher { return Matcher{__a: nil} }`)
		t.Ok(vars["__a"])

		t.End()
	})

	Test(t, "CompareDecl: non-matching decl returns nil", func(t *T) {
		other := parseDecl(`func Other() {}`)
		vars := compare.CompareDecl(other, `func Match() Matcher { return Matcher{__a: nil} }`)
		t.Ok(vars == nil)
		t.End()
	})

	Test(t, "CompareDecl: unparsable pattern returns nil", func(t *T) {
		vars := compare.CompareDecl(decl, `not valid go {{{{`)
		t.Ok(vars == nil)
		t.End()
	})

	Test(t, "CompareDecl: empty body pattern matches empty body decl", func(t *T) {
		empty := parseDecl(`func Match() Matcher { return Matcher{} }`)
		vars := compare.CompareDecl(empty, `func Match() Matcher { return Matcher{} }`)
		t.Ok(vars)

		t.End()
	})

	Test(t, "CompareDecl: nil node returns nil", func(t *T) {
		vars := compare.CompareDecl(nil, `func Match() Matcher { return Matcher{} }`)
		t.Ok(vars == nil)
		t.End()
	})

	Test(t, "CompareDecl: pattern with no decl returns nil", func(t *T) {
		vars := compare.CompareDecl(decl, `// just a comment`)
		t.Ok(vars == nil)
		t.End()
	})

	Test(t, "CompareDecl: doc comment on decl is ignored", func(t *T) {
		f, err := parser.ParseFile(token.NewFileSet(), "", "package p\n\n// Comment above Match.\nfunc Match() Matcher { return Matcher{} }\n", parser.ParseComments)
		if err != nil || len(f.Decls) == 0 {
			t.Ok(false)
			t.End()
			return
		}
		vars := compare.CompareDecl(f.Decls[0], `func Match() Matcher { return Matcher{} }`)
		t.Ok(vars)

		t.End()
	})
}
