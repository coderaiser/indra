package removeunusedvariable

import (
	"go/ast"
	"go/token"

	"coderaiser/indra/internal/engine"
)

var Plugin = engine.Plugin{
	Name:   "remove-unused-variable",
	Report: func() string { return "remove unused variable" },
	Traverse: func() map[string]engine.TraverseVisitor {
		return map[string]engine.TraverseVisitor{
			"*ast.BlockStmt": visitBlock,
		}
	},
}

func visitBlock(node ast.Node, _ engine.Vars) []engine.Place {
	block := node.(*ast.BlockStmt)

	type decl struct {
		name string
	}
	var decls []decl
	seen := map[string]bool{}

	for _, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE {
			continue
		}
		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok || ident.Name == "_" {
				continue
			}
			if !seen[ident.Name] {
				seen[ident.Name] = true
				decls = append(decls, decl{name: ident.Name})
			}
		}
	}

	if len(decls) == 0 {
		return nil
	}

	// count reads: ident uses not on the lhs of a := statement
	reads := map[string]int{}
	for _, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if ok && assign.Tok == token.DEFINE {
			// only count rhs of :=
			for _, rhs := range assign.Rhs {
				countIdents(rhs, reads)
			}
			continue
		}
		countIdents(stmt, reads)
	}

	var places []engine.Place
	for _, d := range decls {
		if reads[d.name] == 0 {
			places = append(places, engine.Place{
				Message: "remove unused variable: " + d.name,
			})
		}
	}
	return places
}

func countIdents(n ast.Node, reads map[string]int) {
	ast.Inspect(n, func(node ast.Node) bool {
		ident, ok := node.(*ast.Ident)
		if ok {
			reads[ident.Name]++
		}
		return true
	})
}
