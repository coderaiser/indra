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
			if hasMismatch(p) {
				push(p)
			}
		},
	}
}

func Fix(p Path, _ map[string]any) { applyFix(p) }

// extractVerb returns the canonical verb for the first t.X method call in the
// callback body.
func extractVerb(fnLit *ast.FuncLit) string {
	for _, stmt := range fnLit.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		id, ok := sel.X.(*ast.Ident)
		if !ok || id.Name != "t" {
			continue
		}
		if sel.Sel.Name == "Report" {
			return "report"
		}
		if sel.Sel.Name == "Transform" {
			return "transform"
		}
		if sel.Sel.Name == "NoReport" {
			return "no report"
		}
		if sel.Sel.Name == "NoTransform" {
			return "no transform"
		}
	}

	return ""
}

// extractFixtureName walks a callback body Path and returns the fixture name
// from the first t.Report/Transform/NoReport/NoTransform call.
func extractFixtureName(bodyPath Path) string {
	result := ""
	bodyPath.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			if result != "" {
				return
			}
			call := callPath.Node.(*ast.CallExpr)
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

// hasMismatch returns true if any Test(t, msg, ...) call has a verb segment
// that does not match the callback method, or a missing/wrong fixture name at
// the end of the message.
func hasMismatch(p Path) bool {
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			call := callPath.Node.(*ast.CallExpr)
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
			verb := extractVerb(fnLit)
			if verb == "" {
				return
			}
			bodyPath := Path{Node: fnLit.Body, Stack: append(callPath.Stack, callPath.Node)}
			fixtureName := extractFixtureName(bodyPath)
			if fixtureName == "" {
				return
			}
			msg := msgLit.Value
			if len(msg) >= 2 {
				inner := msg[1 : len(msg)-1]
				parts := strings.Split(inner, ": ")
				if len(parts) < 2 {
					found = true
					callPath.Stop()
					return
				}
				if parts[len(parts)-1] != fixtureName {
					found = true
					callPath.Stop()
					return
				}
				if len(parts) >= 3 && parts[1] != verb {
					found = true
					callPath.Stop()
				}
			}
		},
	})
	return found
}

// applyFix rewrites each Test message so its verb segment matches the callback
// method and its last ": "-separated segment equals the fixture name.
func applyFix(path Path) {
	path.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(callPath Path) {
			fixTestCall(callPath)
		},
	})
}

func fixTestCall(callPath Path) {
	call, ok := callPath.Node.(*ast.CallExpr)
	if !ok {
		return
	}

	msgLit, fnLit, ok := testCallParts(call)
	if !ok {
		return
	}

	verb := extractVerb(fnLit)
	if verb == "" {
		return
	}

	fixtureName := extractFixtureName(Path{
		Node:  fnLit.Body,
		Stack: append(callPath.Stack, callPath.Node),
	})
	if fixtureName == "" {
		return
	}

	fixMessage(msgLit, verb, fixtureName)
}

func testCallParts(call *ast.CallExpr) (*ast.BasicLit, *ast.FuncLit, bool) {
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "Test" || len(call.Args) < 3 {
		return nil, nil, false
	}

	msgLit, ok := call.Args[1].(*ast.BasicLit)
	if !ok {
		return nil, nil, false
	}

	fnLit, ok := call.Args[2].(*ast.FuncLit)
	if !ok {
		return nil, nil, false
	}

	return msgLit, fnLit, true
}

func fixMessage(msgLit *ast.BasicLit, verb, fixtureName string) {
	msg := msgLit.Value
	if len(msg) < 2 {
		return
	}

	inner := msg[1 : len(msg)-1]
	parts := strings.Split(inner, ": ")

	changed := false

	switch len(parts) {
	case 0, 1:
		parts = []string{verb, fixtureName}
		changed = true

	case 2:
		if parts[1] != fixtureName {
			parts = append(parts, fixtureName)
			changed = true
		}

	default:
		if parts[1] != verb {
			parts[1] = verb
			changed = true
		}

		if parts[len(parts)-1] != fixtureName {
			parts[len(parts)-1] = fixtureName
			changed = true
		}
	}

	if changed {
		msgLit.Value = `"` + strings.Join(parts, ": ") + `"`
	}
}


// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Traverse() Traverser             { return Traverse() }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
