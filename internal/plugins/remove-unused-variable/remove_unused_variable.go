package remove_unused_variable

import (
	"go/ast"
	"go/token"

	. "coderaiser/indra/types"
)

func Report(node ast.Node) string {
	if node == nil {
		return "remove unused variable"
	}
	unused := unusedVarNames(node.(*ast.BlockStmt))
	if len(unused) == 0 {
		return "remove unused variable"
	}
	return "remove unused variable: " + unused[0]
}

func Traverse() Traverser {
	return Traverser{
		"*ast.BlockStmt": findUnusedVars,
	}
}

func findUnusedVars(node ast.Node, push func(ast.Node)) {
	block := node.(*ast.BlockStmt)
	if len(unusedVarNames(block)) > 0 {
		push(block)
	}
}

// unusedVarNames returns the names declared via `:=` in block that are never
// read afterwards.
func unusedVarNames(block *ast.BlockStmt) []string {
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

	var unused []string
	for _, d := range decls {
		if reads[d] == 0 {
			unused = append(unused, d)
		}
	}
	return unused
}

// Fix removes unused variables from a block in place.
// node is *ast.BlockStmt; options is unused.
func Fix(node ast.Node, _ map[string]any) {
	block := node.(*ast.BlockStmt)
	unused := make(map[string]bool, len(unusedVarNames(block)))
	for _, n := range unusedVarNames(block) {
		unused[n] = true
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
