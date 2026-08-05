package compare

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// Vars holds bound pattern holes keyed by hole name.
type Vars = map[string]ast.Node

// ArgSlice is a sentinel node so a CallExpr arg list fits in Vars.
type ArgSlice struct {
	Args []ast.Expr
}

func (ArgSlice) Pos() token.Pos { return token.NoPos }
func (ArgSlice) End() token.Pos { return token.NoPos }

// BodySlice is a sentinel node so a FuncLit body stmt list fits in Vars.
type BodySlice struct {
	Stmts []ast.Stmt
}

func (BodySlice) Pos() token.Pos { return token.NoPos }
func (BodySlice) End() token.Pos { return token.NoPos }

// Compare matches node against a pattern string.
// It returns the bound holes (Vars) or nil when there is no match.
func Compare(node ast.Node, pattern string) Vars {
	pat := parsePattern(pattern)
	if pat == nil {
		return nil
	}
	vars := make(Vars)
	if !matchNode(pat, node, vars) {
		return nil
	}
	return vars
}

// parsePattern parses a statement-level pattern into an ast.Node.
// The __body sentinel is preprocessed so the pattern parses as valid Go.
func parsePattern(s string) ast.Node {
	s = strings.ReplaceAll(s, "{ __body }", "{ __body() }")
	file, err := parser.ParseFile(token.NewFileSet(), "pattern.go", "package p; func _() { "+s+" }", 0)
	if err != nil {
		return nil
	}
	return file.Decls[0].(*ast.FuncDecl).Body.List[0]
}

// matchNode checks pat against real, while recording hole bindings in vars.
func matchNode(pat ast.Node, real ast.Node, vars Vars) bool {
	if pat == nil {
		return real == nil
	}

	// __body sentinel: a block containing only __body() captures the real
	// block's statement list as BodySlice.
	if patBlock, ok := pat.(*ast.BlockStmt); ok && len(patBlock.List) == 1 {
		if expr, ok := patBlock.List[0].(*ast.ExprStmt); ok {
			if call, ok := expr.X.(*ast.CallExpr); ok {
				if fn, ok := call.Fun.(*ast.Ident); ok && fn.Name == "__body" {
					if realBlock, ok := real.(*ast.BlockStmt); ok {
						vars["__body"] = BodySlice{Stmts: realBlock.List}
						return true
					}
				}
			}
		}
	}

	// __ discard hole: match any node, do not store.
	if ident, ok := pat.(*ast.Ident); ok && ident.Name == "__" {
		return true
	}

	// named hole __name: bind any node, enforce linked equality on repeat.
	if ident, ok := pat.(*ast.Ident); ok && strings.HasPrefix(ident.Name, "__") && ident.Name != "__" {
		// __array only matches a CompositeLit whose type is an ArrayType.
		if ident.Name == "__array" {
			lit, ok := real.(*ast.CompositeLit)
			if !ok {
				return false
			}
			if _, ok := lit.Type.(*ast.ArrayType); !ok {
				return false
			}
			return bind("__array", real, vars)
		}
		return bind(ident.Name, real, vars)
	}

	if real == nil {
		return false
	}

	// For BlockStmt patterns (function bodies) the sentinel is handled above;
	// otherwise require matching block structure normally.

	if reflect.TypeOf(pat) != reflect.TypeOf(real) {
		return false
	}

	// idents: matches require equal names (holes are handled above; types
	// already align after the reflect.TypeOf check, so this assertion always
	// succeeds).
	if patIdent, ok := pat.(*ast.Ident); ok {
		realIdent := real.(*ast.Ident)
		return patIdent.Name == realIdent.Name
	}

	return matchChildren(pat, real, vars)
}

// bind records a binding, enforcing that linked holes match the same source.
func bind(name string, real ast.Node, vars Vars) bool {
	if existing, ok := vars[name]; ok {
		return printed(existing) == printed(real)
	}
	vars[name] = real
	return true
}

// printed returns the canonical go/format output for a node.
func printed(n ast.Node) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	_ = format.Node(&buf, token.NewFileSet(), n)
	return buf.String()
}

// matchChildren walks the exported struct fields of pat and real via
// reflect and recursively matches each non-position field.
func matchChildren(pat ast.Node, real ast.Node, vars Vars) bool {
	rp := reflect.ValueOf(pat)
	rr := reflect.ValueOf(real)

	if rp.Kind() != reflect.Pointer || rp.Elem().Kind() != reflect.Struct {
		return true
	}
	if rr.Kind() != reflect.Pointer || rr.Elem().Kind() != reflect.Struct {
		return false
	}

	p := rp.Elem()
	r := rr.Elem()
	pt := p.Type()

	for i := 0; i < p.NumField(); i++ {
		pf := p.Field(i)
		rf := r.Field(i)
		ft := pt.Field(i).Type

		// skip token.Pos bookkeeping fields.
		if ft == reflect.TypeOf(token.NoPos) {
			continue
		}

		// skip comment fields — comments don't affect pattern matching.
		if ft == reflect.TypeOf(&ast.CommentGroup{}) {
			continue
		}

		if !matchField(pt.Field(i).Name, pf, rf, vars) {
			return false
		}
	}

	return true
}

// matchField matches a single struct field pair (pattern vs real).
func matchField(name string, pf, rf reflect.Value, vars Vars) bool {
	switch pf.Kind() {
	case reflect.Slice:
		if pf.IsNil() {
			return rf.IsNil()
		}
		return matchSliceLike(pf, rf, vars)
	case reflect.Pointer:
		if pf.IsNil() {
			return rf.IsNil()
		}
		if rf.IsNil() {
			return false
		}
		return matchNode(pf.Interface().(ast.Node), rf.Interface().(ast.Node), vars)
	case reflect.Interface:
		if pf.IsNil() {
			return rf.IsNil()
		}
		if rf.IsNil() {
			return false
		}
		on, okP := pf.Interface().(ast.Node)
		rn, okR := rf.Interface().(ast.Node)
		if !okP || !okR {
			return pf.Interface() == rf.Interface()
		}
		return matchNode(on, rn, vars)
	default:
		return pf.Interface() == rf.Interface()
	}
}

// matchSliceLike matches a slice field, handling the __args hole.
func matchSliceLike(pf, rf reflect.Value, vars Vars) bool {
	// A slice containing only a __args hole captures the whole real slice.
	if pf.Len() == 1 {
		expr, ok := pf.Index(0).Interface().(ast.Expr)
		if ok {
			if ident, ok := expr.(*ast.Ident); ok && ident.Name == "__args" {
				var args []ast.Expr
				for i := 0; i < rf.Len(); i++ {
					if e, ok := rf.Index(i).Interface().(ast.Expr); ok {
						args = append(args, e)
					}
				}
				vars["__args"] = ArgSlice{Args: args}
				return true
			}
		}
	}

	if pf.Len() != rf.Len() {
		return false
	}

	for i := 0; i < pf.Len(); i++ {
		if pf.Index(i).Kind() == reflect.Pointer || pf.Index(i).Kind() == reflect.Interface {
			if !matchNode(pf.Index(i).Interface().(ast.Node), rf.Index(i).Interface().(ast.Node), vars) {
				return false
			}
		} else if pf.Index(i).Interface() != rf.Index(i).Interface() {
			return false
		}
	}

	return true
}
