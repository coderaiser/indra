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
	if ruleName == "" {
		return
	}
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

// hasMissingPrefix returns true if any Test(t, msg, ...) call in the file has
// a msg not starting with "<ruleName>: ".
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
			if len(msg) >= 2 && !strings.HasPrefix(msg[1:len(msg)-1], prefix) {
				found = true
				cp.Stop()
			}
		},
	})
	return found
}

// applyPrefix prepends "<ruleName>: " to each Test message that lacks it.
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
				if !strings.HasPrefix(inner, prefix) {
					lit.Value = `"` + prefix + inner + `"`
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
