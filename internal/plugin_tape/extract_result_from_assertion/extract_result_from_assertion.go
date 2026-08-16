package extract_result_from_assertion

import (
	"go/ast"

	. "coderaiser/indra/types"
)

func Report() string { return "extract inline expression from assertion" }

// Match guards the call-extraction patterns so a function-call result is not
// re-extracted when a "result" variable is already declared in the containing
// block (which would shadow the injected declaration).
func Match() Matcher {
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

// noResultInBlock is a guard that rejects re-extraction when a "result"
// variable is already declared in the containing block (which would shadow the
// injected declaration). Statement matching always runs inside a block.
func noResultInBlock(_ Vars, block *ast.BlockStmt) bool {
	return !blockDeclares(block, "result")
}

// blockDeclares reports whether any statement in block declares name via a
// short variable declaration (:=).
func blockDeclares(block *ast.BlockStmt, name string) bool {
	for _, s := range block.List {
		if Compare(s, name+" := __a") {
			return true
		}
	}
	return false
}

// Plugin wraps the rule for the registry: a replacer with a Match guard.
type Plugin struct{}

func (Plugin) Report() string    { return Report() }
func (Plugin) Match() Matcher    { return Match() }
func (Plugin) Replace() Replacer { return Replace() }
