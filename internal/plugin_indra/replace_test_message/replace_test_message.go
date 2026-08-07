package replace_test_message

import (
	"go/ast"
	"strings"

	. "coderaiser/indra/types"
)

// fixtureMethods are the callback methods whose first string argument names a
// fixture file used by the surrounding Test(t, msg, ...) call.
var fixtureMethods = map[string]bool{
	"Report": true, "Transform": true,
	"NoReport": true, "NoTransform": true,
}

func Report(_ Path) string { return "replace test message" }

func Traverse() Traverser {
	return Traverser{
		"*ast.File": func(p Path, push func(Path)) {
			if hasMissingFixtureName(p) {
				push(p)
			}
		},
	}
}

func Fix(p Path, _ map[string]any) { applyFixtureNames(p) }

// afterSeparator returns the part of inner that follows the first ": " rule
// prefix separator. When there is no separator the whole string is returned.
func afterSeparator(inner string) string {
	idx := strings.Index(inner, ": ")
	if idx < 0 {
		return inner
	}
	return inner[idx+2:]
}

// extractFixtureName walks a callback body Path and returns the fixture name
// from the first t.Report/Transform/NoReport/NoTransform call.
func extractFixtureName(bodyPath Path) string {
	result := ""
	bodyPath.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(cp Path) {
			if result != "" {
				return
			}
			call := cp.Node.(*ast.CallExpr)
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || id.Name != "t" || !fixtureMethods[sel.Sel.Name] {
				return
			}
			if len(call.Args) < 1 {
				return
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				return
			}
			s := lit.Value
			if len(s) >= 2 {
				result = s[1 : len(s)-1]
			}
		},
	})
	return result
}

// hasMissingFixtureName returns true if any Test(t, msg, ...) call whose
// callback names a fixture has a msg that lacks the fixture name after the
// ": " rule separator.
func hasMissingFixtureName(p Path) bool {
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(cp Path) {
			if found {
				return
			}
			call := cp.Node.(*ast.CallExpr)
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "Test" || len(call.Args) < 3 {
				return
			}
			msgLit, ok := call.Args[1].(*ast.BasicLit)
			if !ok {
				return
			}
			fnLit, ok := call.Args[2].(*ast.FuncLit)
			if !ok {
				return
			}
			bodyPath := Path{Node: fnLit.Body, Stack: append(cp.Stack, cp.Node)}
			fixtureName := extractFixtureName(bodyPath)
			if fixtureName == "" {
				return
			}
			msg := msgLit.Value
			if len(msg) >= 2 && !strings.Contains(afterSeparator(msg[1:len(msg)-1]), fixtureName) {
				found = true
			}
		},
	})
	return found
}

// applyFixtureNames appends " <fixtureName>" to every Test message that is
// missing its fixture name after the ": " separator.
func applyFixtureNames(p Path) {
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(cp Path) {
			call := cp.Node.(*ast.CallExpr)
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "Test" || len(call.Args) < 3 {
				return
			}
			msgLit, ok := call.Args[1].(*ast.BasicLit)
			if !ok {
				return
			}
			fnLit, ok := call.Args[2].(*ast.FuncLit)
			if !ok {
				return
			}
			bodyPath := Path{Node: fnLit.Body, Stack: append(cp.Stack, cp.Node)}
			fixtureName := extractFixtureName(bodyPath)
			if fixtureName == "" {
				return
			}
			msg := msgLit.Value
			if len(msg) >= 2 {
				inner := msg[1 : len(msg)-1]
				if !strings.Contains(afterSeparator(inner), fixtureName) {
					msgLit.Value = `"` + inner + " " + fixtureName + `"`
				}
			}
		},
	})
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
