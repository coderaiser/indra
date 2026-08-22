package apply_exports_order

import (
	"go/ast"
	"slices"

	. "coderaiser/indra/types"
)

// shapes holds the canonical export orders enforced by this rule:
//   - Replacer:  Report, Match, Replace
//   - Traverser: Report, Fix, Traverse
//   - Includer:  Report, Include, Fix, Filter
//
// A file whose top-level exported functions contain every name of a shape is
// checked against that shape's order.
var shapes = [][]string{
	{"Report", "Match", "Replace"},
	{"Report", "Fix", "Traverse"},
	{"Report", "Include", "Fix", "Filter"},
}

func Report(_ Path) string { return "Apply exports order" }

// Fix reorders the top-level exported functions of file into the detected
// shape's canonical order. Non-function declarations and unexported functions
// keep their positions; exported functions outside the shape follow the shape
// functions preserving their relative order. node is *ast.File; options is
// unused.
func Fix(p Path, _ map[string]any) {
	file := p.Node.(*ast.File)
	fns := topLevelFuncs(file)
	shape := detectShape(funcNames(fns))
	if shape == nil {
		return
	}
	file.Decls = rebuildDecls(file.Decls, reorder(fns, shape))
}

func Traverse() Traverser {
	return Traverser{
		"*ast.File": func(p Path, push func(Path)) {
			file := p.Node.(*ast.File)
			names := funcNames(topLevelFuncs(file))
			shape := detectShape(names)
			if shape == nil {
				return
			}
			if !inOrder(names, shape) {
				push(p)
			}
		},
	}
}

// topLevelFuncs returns the exported top-level function declarations of file
// in declaration order.
func topLevelFuncs(file *ast.File) []*ast.FuncDecl {
	var fns []*ast.FuncDecl
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && ast.IsExported(fn.Name.Name) {
			fns = append(fns, fn)
		}
	}
	return fns
}

// detectShape returns the first shape whose every name occurs in names, or nil
// when the file does not implement any known plugin shape.
func detectShape(names []string) []string {
	for _, shape := range shapes {
		if containsAll(names, shape) {
			return shape
		}
	}
	return nil
}

func containsAll(names, shape []string) bool {
	for _, s := range shape {
		if !slices.Contains(names, s) {
			return false
		}
	}
	return true
}

// inOrder reports whether every shape name appears in names at a non-decreasing
// index, i.e. the declarations already follow the canonical order.
func inOrder(names, shape []string) bool {
	prev := -1
	for _, s := range shape {
		i := slices.Index(names, s)
		if i < prev {
			return false
		}
		prev = i
	}
	return true
}

// reorder returns fns sorted by shape first, followed by the remaining
// functions in their original relative order.
func reorder(fns []*ast.FuncDecl, shape []string) []*ast.FuncDecl {
	byName := make(map[string]*ast.FuncDecl, len(fns))
	for _, fn := range fns {
		byName[fn.Name.Name] = fn
	}
	result := make([]*ast.FuncDecl, 0, len(fns))
	for _, name := range shape {
		if fn, ok := byName[name]; ok {
			result = append(result, fn)
			delete(byName, name)
		}
	}
	for _, fn := range fns {
		if _, ok := byName[fn.Name.Name]; ok {
			result = append(result, fn)
		}
	}
	return result
}

// rebuildDecls replaces every exported top-level function declaration slot in
// decls with the next function from fns, keeping all other declarations in
// place.
func rebuildDecls(decls []ast.Decl, fns []*ast.FuncDecl) []ast.Decl {
	result := make([]ast.Decl, 0, len(decls))
	fi := 0
	for _, decl := range decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Recv == nil && ast.IsExported(fn.Name.Name) {
			if fi < len(fns) {
				result = append(result, fns[fi])
				fi++
			}
		} else {
			result = append(result, decl)
		}
	}
	return result
}

func funcNames(fns []*ast.FuncDecl) []string {
	names := make([]string, len(fns))
	for i, fn := range fns {
		names[i] = fn.Name.Name
	}
	return names
}

// Plugin wraps the rule for the registry: an AST-walking plugin.
type Plugin struct{}

func (Plugin) Report(p Path) string            { return Report(p) }
func (Plugin) Fix(p Path, opts map[string]any) { Fix(p, opts) }
func (Plugin) Traverse() Traverser             { return Traverse() }
