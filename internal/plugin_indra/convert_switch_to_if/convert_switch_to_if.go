package convert_switch_to_if

import (
	"go/ast"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/ast/astutil"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "use 'if' instead of 'switch'" }

func Traverse() Traverser {
	return Traverser{
		"*ast.SwitchStmt": func(p Path, push func(Path)) {
			if isConvertible(p) {
				push(p)
			}
		},
	}
}

// isConvertible reports whether the switch can become a chain of ifs: a
// tag-based switch (no Init) where every non-default case has exactly one
// value and contains a return (and no fallthrough).
func isConvertible(p Path) bool {
	sw := p.Node.(*ast.SwitchStmt)
	if sw.Tag == nil || sw.Init != nil {
		return false
	}
	for _, clause := range sw.Body.List {
		cc := clause.(*ast.CaseClause)
		if cc.List == nil || len(cc.List) != 1 {
			return false
		}
		caseBodyPath := Path{Node: &ast.BlockStmt{List: cc.Body}, Stack: p.Stack}
		if hasFallthrough(caseBodyPath) || !hasReturn(caseBodyPath) {
			return false
		}
	}
	return true
}

// Fix replaces the switch with a chain of if tag == val { ... } statements.
// The shared sw.Tag node is reused as X in every generated BinaryExpr.
func Fix(p Path, _ map[string]any) {
	sw := p.Node.(*ast.SwitchStmt)
	for _, clause := range sw.Body.List {
		cc := clause.(*ast.CaseClause)
		p.InsertBefore(&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  sw.Tag,
				Op: token.EQL,
				Y:  cc.List[0],
			},
			Body: &ast.BlockStmt{List: removeBreaks(cc.Body)},
		})
	}
	p.Delete()
	// Normalize the enclosing block's source positions so go/printer flows the
	// new ifs and the trailing statements contiguously instead of "catching up"
	// to the stale positions left behind by the removed switch.
	if block := switchBlock(p); block != nil {
		stripPositions(block)
	}
}

// switchBlock returns the *ast.BlockStmt that directly contains the switch, or
// nil when the switch has no parent or its parent is not a block.
func switchBlock(p Path) *ast.BlockStmt {
	parent, ok := p.ParentPath()
	if !ok {
		return nil
	}
	block, ok := parent.Node.(*ast.BlockStmt)
	if !ok {
		return nil
	}
	return block
}

// stripPositions zeroes every token.Pos field in the sub-tree rooted at root.
// Generated nodes that reuse original AST nodes would otherwise carry stale
// source positions that make go/printer insert spurious line breaks.
func stripPositions(root ast.Node) {
	posType := reflect.TypeOf(token.Pos(0))
	astutil.Apply(root, func(c *astutil.Cursor) bool {
		if c.Node() == nil {
			return true
		}
		v := reflect.ValueOf(c.Node()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Type() == posType {
				v.Field(i).SetInt(0)
			}
		}
		return true
	}, nil)
}

// hasReturn reports whether the sub-tree rooted at p contains a ReturnStmt.
func hasReturn(p Path) bool {
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.ReturnStmt": func(_ Path) { found = true },
	})
	return found
}

// hasFallthrough reports whether the sub-tree rooted at p contains a
// fallthrough BranchStmt.
func hasFallthrough(p Path) bool {
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.BranchStmt": func(bp Path) {
			if bp.Node.(*ast.BranchStmt).Tok == token.FALLTHROUGH {
				found = true
			}
		},
	})
	return found
}

// removeBreaks drops break statements from a case body. A break at the end of a
// case body is meaningless once the switch becomes an if.
func removeBreaks(stmts []ast.Stmt) []ast.Stmt {
	out := make([]ast.Stmt, 0, len(stmts))
	for _, s := range stmts {
		if b, ok := s.(*ast.BranchStmt); ok && b.Tok == token.BREAK {
			continue
		}
		out = append(out, s)
	}
	return out
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
