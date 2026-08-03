package compare_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"coderaiser/indra/compare"
	tape "github.com/coderaiser/go-tape"
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

	tape.Test(t, "CompareDecl: matching decl returns non-nil vars", func(t *tape.T) {
		vars := compare.CompareDecl(decl, `func Match() Matcher { return Matcher{__a: nil} }`)
		t.Ok(vars != nil)
		t.End()
	})

	tape.Test(t, "CompareDecl: matching decl binds hole", func(t *tape.T) {
		vars := compare.CompareDecl(decl, `func Match() Matcher { return Matcher{__a: nil} }`)
		t.Ok(vars["__a"] != nil)
		t.End()
	})

	tape.Test(t, "CompareDecl: non-matching decl returns nil", func(t *tape.T) {
		other := parseDecl(`func Other() {}`)
		vars := compare.CompareDecl(other, `func Match() Matcher { return Matcher{__a: nil} }`)
		t.Equal(vars, compare.Vars(nil))
		t.End()
	})

	tape.Test(t, "CompareDecl: unparseable pattern returns nil", func(t *tape.T) {
		vars := compare.CompareDecl(decl, `not valid go {{{{`)
		t.Equal(vars, compare.Vars(nil))
		t.End()
	})

	tape.Test(t, "CompareDecl: empty body pattern matches empty body decl", func(t *tape.T) {
		empty := parseDecl(`func Match() Matcher { return Matcher{} }`)
		vars := compare.CompareDecl(empty, `func Match() Matcher { return Matcher{} }`)
		t.Ok(vars != nil)
		t.End()
	})

	tape.Test(t, "CompareDecl: nil node returns nil", func(t *tape.T) {
		vars := compare.CompareDecl(nil, `func Match() Matcher { return Matcher{} }`)
		t.Equal(vars, compare.Vars(nil))
		t.End()
	})

	tape.Test(t, "CompareDecl: pattern with no decl returns nil", func(t *tape.T) {
		vars := compare.CompareDecl(decl, `// just a comment`)
		t.Equal(vars, compare.Vars(nil))
		t.End()
	})
}
