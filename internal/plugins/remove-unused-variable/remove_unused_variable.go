package remove_unused_variable

import (
	"go/ast"
	"go/token"
	"strings"

	. "coderaiser/indra/types"
)

func Report() string { return "remove unused variable" }

func Traverse() Traverser {
	return Traverser{
		"*ast.BlockStmt": visitBlock,
	}
}

func visitBlock(node ast.Node, _ Vars) []Place {
	block := node.(*ast.BlockStmt)

	var decls []string
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
				decls = append(decls, ident.Name)
			}
		}
	}

	if len(decls) == 0 {
		return nil
	}

	reads := map[string]int{}
	for _, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if ok && assign.Tok == token.DEFINE {
			for _, rhs := range assign.Rhs {
				countIdents(rhs, reads)
			}
			continue
		}
		countIdents(stmt, reads)
	}

	var places []Place
	for _, d := range decls {
		if reads[d] == 0 {
			places = append(places, Place{
				Message: "remove unused variable: " + d,
			})
		}
	}
	return places
}

// Fix removes unused variables from a block in place.
// node is *ast.BlockStmt. places contains findings from Traverse (one per var).
func Fix(node ast.Node, places []Place) {
	block := node.(*ast.BlockStmt)
	unused := make(map[string]bool, len(places))
	for _, p := range places {
		unused[strings.TrimPrefix(p.Message, "remove unused variable: ")] = true
	}
	kept := block.List[:0]
	for _, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE {
			kept = append(kept, stmt)
			continue
		}
		allUnused := true
		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok || ident.Name == "_" {
				continue
			}
			if !unused[ident.Name] {
				allUnused = false
			}
		}
		if allUnused {
			continue // drop statement
		}
		for _, lhs := range assign.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if ok && unused[ident.Name] {
				ident.Name = "_"
			}
		}
		kept = append(kept, stmt)
	}
	block.List = kept
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
