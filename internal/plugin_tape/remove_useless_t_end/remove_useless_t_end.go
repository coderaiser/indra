package remove_useless_t_end

import (
	. "coderaiser/indra/operator"
	. "coderaiser/indra/types"
)

// Report is called with the enclosing Test statement path.
func Report(_ Path) string { return `Avoid useless "t.End()"` }

// Fix is a no-op: duplicate End statements are removed during Traverse, while
// their cursors are live. The engine still calls Fix for each pushed finding.
func Fix(_ Path, _ map[string]any) {
	_ = "removed during traverse"
}

// Traverse finds Test blocks whose body contains more than one t.End():
// every End beyond the first is a runtime no-op. Duplicates are removed in
// place during the sub-walk, while their astutil cursor is still current.
func Traverse() Traverser {
	return Traverser{
		`Test(__a, __b, func(__c *T) { __body })`: func(path Path, push func(Path)) {
			seen := false
			path.Traverse(map[string]func(Path){
				`__x.End()`: func(p Path) {
					if seen {
						Remove(p)
						push(path)
					}
					seen = true
				},
			})
		},
	}
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
