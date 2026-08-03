package engine_runner

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	loader "coderaiser/indra/engine-loader"
	"coderaiser/indra/types"
)

func parse(t *testing.T, src string) (*ast.File, *token.FileSet) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	return file, fset
}

func items(plugins []loader.PluginFuncs) []PluginItem {
	kinds := loader.Load(plugins, loader.Config{})
	out := make([]PluginItem, len(kinds))
	for i, k := range kinds {
		out[i] = PluginItem{Rule: k.Name(), Plugin: k}
	}
	return out
}

func replacerItem() []loader.PluginFuncs {
	return []loader.PluginFuncs{{
		Name:    "eq",
		Report:  func() string { return "use DeepEqual" },
		Match:   func() types.Matcher { return types.Matcher{"t.Equal(__a, __b)": nil} },
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"} },
	}}
}

func traverserFuncs(name, key string) loader.PluginFuncs {
	return loader.PluginFuncs{
		Name:   name,
		Report: func() string { return "issue" },
		Traverse: func() types.Traverser {
			return types.Traverser{
				key: func(node ast.Node, vars types.Vars) []types.Place {
					return []types.Place{{Message: "issue"}}
				},
			}
		},
		Fix: func(node ast.Node, places []types.Place) {},
	}
}

func printFile(t *testing.T, file *ast.File, fset *token.FileSet) string {
	t.Helper()
	if fset == nil {
		fset = token.NewFileSet()
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		t.Fatalf("format: %v", err)
	}
	return buf.String()
}

func TestRunReplacerReturnsPlaces(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(places))
	}
	if places[0].Rule != "eq" || places[0].Message != "use DeepEqual" {
		t.Fatalf("unexpected place: %+v", places[0])
	}
}

func TestRunReplacerFixRewrites(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)
	if !strings.Contains(out, "t.DeepEqual(a, b)") {
		t.Fatalf("expected DeepEqual after fix:\n%s", out)
	}
}

func TestRunMsgOverridesReport(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	kinds := loader.Load(replacerItem(), loader.Config{})
	pl := []PluginItem{{Rule: "eq", Plugin: kinds[0], Msg: "custom"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if places[0].Message != "custom" {
		t.Fatalf("expected custom message, got %q", places[0].Message)
	}
}

func TestRunTraverserReturnsPlaces(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	pl := items([]loader.PluginFuncs{traverserFuncs("file", "*ast.File")})
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 1 {
		t.Fatalf("expected 1 place from traverser, got %d", len(places))
	}
}

func TestRunTraverserFixCallsFix(t *testing.T) {
	src := "package p\n\nfunc f() {}\n"
	file, fset := parse(t, src)
	funcs := loader.PluginFuncs{
		Name:   "file",
		Report: func() string { return "file issue" },
		Traverse: func() types.Traverser {
			return types.Traverser{
				"*ast.File": func(node ast.Node, vars types.Vars) []types.Place {
					return []types.Place{{Message: "file issue"}}
				},
			}
		},
		Fix: func(node ast.Node, places []types.Place) {
			f := node.(*ast.File)
			f.Name.Name = "q"
		},
	}
	pl := items([]loader.PluginFuncs{funcs})
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	if file.Name.Name != "q" {
		t.Fatalf("expected Fix to rename package to q, got %q", file.Name.Name)
	}
}

func TestRunTraverserBlockVisitor(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	pl := items([]loader.PluginFuncs{traverserFuncs("block", "*ast.BlockStmt")})
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 1 {
		t.Fatalf("expected 1 place from block traverser, got %d", len(places))
	}
}

func TestRunEmptyPlugins(t *testing.T) {
	src := "package p\n\nfunc f() {}\n"
	file, fset := parse(t, src)
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: nil})
	if len(places) != 0 {
		t.Fatalf("expected no places with empty plugins, got %d", len(places))
	}
}

func TestRunGuardRejects(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	funcs := loader.PluginFuncs{
		Name:   "guarded",
		Report: func() string { return "m" },
		Match: func() types.Matcher {
			return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return false }}
		},
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "x"} },
	}
	pl := items([]loader.PluginFuncs{funcs})
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 0 {
		t.Fatalf("expected no places when guard rejects, got %d", len(places))
	}
}

func TestRunTraverserBlockFix(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	funcs := loader.PluginFuncs{
		Name:   "block",
		Report: func() string { return "block" },
		Traverse: func() types.Traverser {
			return types.Traverser{
				"*ast.BlockStmt": func(node ast.Node, vars types.Vars) []types.Place {
					return []types.Place{{Message: "block"}}
				},
			}
		},
		Fix: func(node ast.Node, places []types.Place) {
			block := node.(*ast.BlockStmt)
			block.List = nil
		},
	}
	pl := items([]loader.PluginFuncs{funcs})
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	fileBlock := file.Decls[0].(*ast.FuncDecl).Body
	if len(fileBlock.List) != 0 {
		t.Fatalf("expected block fix to empty list, got %d", len(fileBlock.List))
	}
}

func TestRunTraverserMsgOverride(t *testing.T) {
	src := "package p\n\nfunc f() {}\n"
	file, fset := parse(t, src)
	kinds := loader.Load([]loader.PluginFuncs{traverserFuncs("file", "*ast.File")}, loader.Config{})
	pl := []PluginItem{{Rule: "file", Plugin: kinds[0], Msg: "override"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if places[0].Message != "override" {
		t.Fatalf("expected traverser msg override, got %q", places[0].Message)
	}
}

func TestRunFixCountLoops(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n\tt.Equal(c, d)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 3, Plugins: pl})
	out := printFile(t, file, fset)
	if strings.Count(out, "DeepEqual") != 2 {
		t.Fatalf("expected both statements converted via loop:\n%s", out)
	}
}

func TestRunNoFixNoRewrite(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: false, Plugins: pl})
	out := printFile(t, file, fset)
	if strings.Contains(out, "DeepEqual") {
		t.Fatalf("expected no rewrite without fix:\n%s", out)
	}
}

func TestSubstituteUnmatchedHole(t *testing.T) {
	got := substitute("hello __nope", Vars{})
	if got != "hello __nope" {
		t.Fatalf("expected hole preserved, got %q", got)
	}
}

func TestSubstituteMatchedHole(t *testing.T) {
	src := "package p\n\nfunc f() { x := 1 }\n"
	file, _ := parse(t, src)
	stmt := file.Decls[0].(*ast.FuncDecl).Body.List[0]
	got := substitute("a = __x", Vars{"__x": stmt})
	if !strings.Contains(got, "x := 1") {
		t.Fatalf("expected substituted source, got %q", got)
	}
}

func TestRunRenderArgSlice(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tuse(f, g)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "args",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{"use(__args)": nil} },
		Replace: func() types.Replacer { return types.Replacer{"use(__args)": "wrap(__args)"} },
	}}
	pl := items(funcs)
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)
	if !strings.Contains(out, "wrap(f, g)") {
		t.Fatalf("expected wrapped arg slice:\n%s", out)
	}
}

func TestRunRenderBodySlice(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tTest(\"s\", \"n\", func(t *T) {\n\t\tt.Equal(a, b)\n\t})\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:   "body",
		Report: func() string { return "m" },
		Match:  func() types.Matcher { return types.Matcher{"Test(__a, __b, func(__r *T) { __body })": nil} },
		Replace: func() types.Replacer {
			return types.Replacer{"Test(__a, __b, func(__r *T) { __body })": "Test(__a, __b, func(__r *T) {\n__body\n__r.End()\n})"}
		},
	}}
	pl := items(funcs)
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)
	if !strings.Contains(out, ".End()") {
		t.Fatalf("expected End() added to body:\n%s", out)
	}
}

func TestPrintNodeNil(t *testing.T) {
	if got := printNode(nil); got != "" {
		t.Fatalf("expected empty string for nil node, got %q", got)
	}
}

func TestApplyRewritesMultipleInBlock(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tmakeSlices(x)\n\tmakeSlices(y)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "multi",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{"makeSlices(__x)": nil} },
		Replace: func() types.Replacer { return types.Replacer{"makeSlices(__x)": "v := __x"} },
	}}
	pl := items(funcs)
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)
	if strings.Count(out, "v := ") != 2 {
		t.Fatalf("expected two rewrites in block:\n%s", out)
	}
}

func TestRunUnparseableReplace(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "badreplace",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{"t.Equal(__a, __b)": nil} },
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "func ("} },
	}}
	pl := items(funcs)
	places := RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	if len(places) != 1 {
		t.Fatalf("expected 1 place even with unparseable replace, got %d", len(places))
	}
}

func TestRunTraverserBlockMsgOverride(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	kinds := loader.Load([]loader.PluginFuncs{traverserFuncs("block", "*ast.BlockStmt")}, loader.Config{})
	pl := []PluginItem{{Rule: "block", Plugin: kinds[0], Msg: "blockmsg"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if places[0].Message != "blockmsg" {
		t.Fatalf("expected block traverser msg override, got %q", places[0].Message)
	}
}

func TestSubstituteAndParseError(t *testing.T) {
	if stmts := substituteAndParse("func (", Vars{}); stmts != nil {
		t.Fatalf("expected nil for unparseable template, got %v", stmts)
	}
}

// declFuncs builds a replacer plugin whose Match and Replace carry a
// top-level declaration pattern (used for decl-level rewrites).
func declFuncs() []loader.PluginFuncs {
	return []loader.PluginFuncs{{
		Name:    "decl",
		Report:  func() string { return "decl issue" },
		Match:   func() types.Matcher { return types.Matcher{`func Match() Matcher { return Matcher{__a: nil} }`: nil} },
		Replace: func() types.Replacer { return types.Replacer{`func Match() Matcher { return Matcher{__a: nil} }`: ""} },
	}}
}

func TestRunDeclRewritesReportsPlace(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	pl := items(declFuncs())
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 1 {
		t.Fatalf("expected 1 place from decl rewrite, got %d", len(places))
	}
}

func TestRunDeclRewritesGuardRejects(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:   "decl",
		Report: func() string { return "decl issue" },
		Match: func() types.Matcher {
			return types.Matcher{`func Match() Matcher { return Matcher{__a: nil} }`: func(v types.Vars) bool { return false }}
		},
		Replace: func() types.Replacer { return types.Replacer{`func Match() Matcher { return Matcher{__a: nil} }`: ""} },
	}}
	pl := items(funcs)
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 0 {
		t.Fatalf("expected no places when decl guard rejects, got %d", len(places))
	}
}

func TestRunDeclRewritesRemovesDecl(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	pl := items(declFuncs())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)
	if strings.Contains(out, "func Match") {
		t.Fatalf("expected declaration removed after fix:\n%s", out)
	}
}

func TestRunDeclRewritesNoMatch(t *testing.T) {
	src := "package p\n\nfunc Other() {}\n"
	file, fset := parse(t, src)
	pl := items(declFuncs())
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if len(places) != 0 {
		t.Fatalf("expected no places when decl does not match, got %d", len(places))
	}
}

func TestRunDeclRewritesMsgOverride(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	kinds := loader.Load(declFuncs(), loader.Config{})
	pl := []PluginItem{{Rule: "decl", Plugin: kinds[0], Msg: "custom decl"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
	if places[0].Message != "custom decl" {
		t.Fatalf("expected decl msg override, got %q", places[0].Message)
	}
}

func TestApplyDeclRewritesKeepsNonEmptyTmpl(t *testing.T) {
	file, _ := parse(t, "package p\n\nfunc Match() Matcher { return Matcher{} }\n")
	applyDeclRewrites(file, []declRewrite{{idx: 0, tmpl: "func Other() {}"}})
	if len(file.Decls) != 1 {
		t.Fatalf("expected declaration kept for non-empty tmpl, got %d decls", len(file.Decls))
	}
}

func TestApplyDeclRewritesRemovesMultiple(t *testing.T) {
	file, _ := parse(t, "package p\n\nfunc A() {}\n\nfunc B() {}\n")
	applyDeclRewrites(file, []declRewrite{{idx: 0, tmpl: ""}, {idx: 1, tmpl: ""}})
	if len(file.Decls) != 0 {
		t.Fatalf("expected both declarations removed, got %d decls", len(file.Decls))
	}
}
