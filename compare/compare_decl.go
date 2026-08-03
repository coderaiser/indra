package compare

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// CompareDecl matches a top-level ast.Decl against a declaration pattern string.
// Returns bound hole vars or nil on no match.
func CompareDecl(node ast.Node, pattern string) Vars {
	if node == nil {
		return nil
	}
	pat := parseDeclPattern(pattern)
	if pat == nil {
		return nil
	}
	vars := make(Vars)
	if !matchNode(pat, node, vars) {
		return nil
	}
	return vars
}

// parseDeclPattern parses a top-level declaration pattern string into an ast.Node.
// The __body sentinel is preprocessed identically to parsePattern.
func parseDeclPattern(s string) ast.Node {
	s = strings.ReplaceAll(s, "{ __body }", "{ __body() }")
	file, err := parser.ParseFile(token.NewFileSet(), "pattern.go", "package p\n"+s, 0)
	if err != nil {
		return nil
	}
	if len(file.Decls) == 0 {
		return nil
	}
	return file.Decls[0]
}
