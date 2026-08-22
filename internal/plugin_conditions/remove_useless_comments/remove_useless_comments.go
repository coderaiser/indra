// Package remove_useless_comments removes separator banner comments made of
// repeated "─" (U+2500) characters, e.g.:
//

package remove_useless_comments

import (
	"go/ast"
	"strings"

	. "coderaiser/indra/types"
)

func Report(_ Path) string { return "remove useless comments" }

// Fix drops every comment node that looks like a banner separator. Comment
// groups that still carry non-banner comments are kept with the banner lines
// filtered out.
func Fix(p Path, _ map[string]any) {
	file := p.Node.(*ast.File)
	var kept []*ast.CommentGroup
	for _, cg := range file.Comments {
		filtered := make([]*ast.Comment, 0, len(cg.List))
		for _, c := range cg.List {
			if strings.Count(c.Text, "─") < 2 {
				filtered = append(filtered, c)
			}
		}
		if len(filtered) > 0 {
			cg.List = filtered
			kept = append(kept, cg)
		}
	}
	file.Comments = kept
}

func Traverse() Traverser {
	return Traverser{
		"*ast.File": func(p Path, push func(Path)) {
			file := p.Node.(*ast.File)
			for _, cg := range file.Comments {
				for _, c := range cg.List {
					if strings.Count(c.Text, "─") >= 2 {
						push(p)
						return
					}
				}
			}
		},
	}
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
