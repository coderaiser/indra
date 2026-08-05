package extract_result_from_assertion

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/types"
)

func Report() string { return "extract inline expression from assertion" }

// Match guards the call-extraction patterns so a function-call result is not
// re-extracted when a "result" variable is already declared in the containing
// block (which would shadow the injected declaration).
func Match() Matcher {
	noResultInBlock := func(vars Vars) bool {
		block, ok := vars["$block"].(*ast.BlockStmt)
		if !ok {
			return true
		}
		return !blockDeclares(block, "result")
	}
	return Matcher{
		"__a.Equal(__b(__args), __c)":     noResultInBlock,
		"__a.DeepEqual(__b(__args), __c)": noResultInBlock,
	}
}

func Replace() Replacer {
	return Replacer{
		"__a.Equal(__b(__args), __c)":     "result := __b(__args)\n__a.Equal(result, __c)",
		"__a.DeepEqual(__b(__args), __c)": "result := __b(__args)\n__a.DeepEqual(result, __c)",
		"__a.Equal(__b, __array)":         "expected := __array\n__a.Equal(__b, expected)",
		"__a.DeepEqual(__b, __array)":     "expected := __array\n__a.DeepEqual(__b, expected)",
	}
}

// blockDeclares reports whether any statement in block declares name via a
// short variable declaration (:=).
func blockDeclares(block *ast.BlockStmt, name string) bool {
	for _, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE {
			continue
		}
		for _, lhs := range assign.Lhs {
			if id, ok := lhs.(*ast.Ident); ok && id.Name == name {
				return true
			}
		}
	}
	return false
}
