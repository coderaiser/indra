package apply_fixture_name_to_message

import (
	"go/ast"
	"strings"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "apply fixture name to message" }

func Traverse() Traverser {
	return Traverser{
		"*ast.File": func(p Path, push func(Path)) {
			ruleName := extractRuleName(p)
			if ruleName == "" {
				return
			}
			if hasMissingPrefix(p, ruleName) {
				push(p)
			}
		},
	}
}

func Fix(p Path, _ map[string]any) {
	ruleName := extractRuleName(p)
	applyPrefix(p, ruleName)
}

// extractRuleName walks the file's declarations using Path.Traverse, finds
// var Test = CreateTest("rule-name", ...) and returns the rule name.
func extractRuleName(p Path) string {
	result := ""
	p.Traverse(map[string]func(Path){
		"*ast.ValueSpec": func(vp Path) {
			vs := vp.Node.(*ast.ValueSpec)
			for i, name := range vs.Names {
				if name.Name != "Test" || i >= len(vs.Values) {
					continue
				}
				call, ok := vs.Values[i].(*ast.CallExpr)
				if !ok {
					continue
				}
				fn, ok := call.Fun.(*ast.Ident)
				if !ok || fn.Name != "CreateTest" || len(call.Args) < 1 {
					continue
				}
				lit, ok := call.Args[0].(*ast.BasicLit)
				if !ok {
					continue
				}
				s := lit.Value
				if len(s) >= 2 {
					result = s[1 : len(s)-1]
				}
			}
		},
	})
	return result
}

var fixtureMethods = map[string]bool{
	"Report": true, "Transform": true,
	"NoReport": true, "NoTransform": true,
}

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

func extractFixtureNameFromTest(cp Path) string {
	call := cp.Node.(*ast.CallExpr)
	if len(call.Args) < 3 {
		return ""
	}
	fnLit, ok := call.Args[2].(*ast.FuncLit)
	if !ok {
		return ""
	}
	bodyPath := Path{Node: fnLit.Body, Stack: append(cp.Stack, cp.Node)}
	return extractFixtureName(bodyPath)
}

// hasMissingPrefix returns true if any Test(t, msg, ...) call in the file has
// a msg not starting with "<ruleName>: " or not ending with ": <fixtureName>".
func hasMissingPrefix(p Path, ruleName string) bool {
	prefix := ruleName + ": "
	found := false
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(cp Path) {
			call := cp.Node.(*ast.CallExpr)
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "Test" || len(call.Args) < 2 {
				return
			}
			lit, ok := call.Args[1].(*ast.BasicLit)
			if !ok {
				return
			}
			msg := lit.Value
			if len(msg) >= 2 {
				inner := msg[1 : len(msg)-1]
				if !strings.HasPrefix(inner, prefix) {
					found = true
					cp.Stop()
					return
				}
				fixtureName := extractFixtureNameFromTest(cp)
				if fixtureName != "" && !strings.HasSuffix(inner, ": "+fixtureName) {
					found = true
					cp.Stop()
				}
			}
		},
	})
	return found
}

// applyPrefix prepends "<ruleName>: " to each Test message that lacks it, and
// replaces the last ": "-separated segment with the fixture name when it is
// missing or wrong.
func applyPrefix(p Path, ruleName string) {
	prefix := ruleName + ": "
	p.Traverse(map[string]func(Path){
		"*ast.CallExpr": func(cp Path) {
			call := cp.Node.(*ast.CallExpr)
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "Test" || len(call.Args) < 2 {
				return
			}
			lit, ok := call.Args[1].(*ast.BasicLit)
			if !ok {
				return
			}
			msg := lit.Value
			if len(msg) >= 2 {
				inner := msg[1 : len(msg)-1]
				changed := false
				if !strings.HasPrefix(inner, prefix) {
					inner = prefix + inner
					changed = true
				}
				fixtureName := extractFixtureNameFromTest(cp)
				if fixtureName != "" && !strings.HasSuffix(inner, ": "+fixtureName) {
					if changed {
						inner = inner + ": " + fixtureName
					} else {
						// inner starts with prefix, and prefix ends with ": ",
						// so there is always at least one segment to re-write.
						parts := strings.Split(inner, ": ")
						parts[len(parts)-1] = fixtureName
						inner = strings.Join(parts, ": ")
					}
					changed = true
				}
				if changed {
					lit.Value = `"` + inner + `"`
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
