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
	. "github.com/coderaiser/go-tape"
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
		Name:   "eq",
		Report: func() string { return "use DeepEqual" },
		Match: func() types.Matcher {
			return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}
		},
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"} },
	}}
}

func traverserFuncs(name, key string) loader.PluginFuncs {
	return loader.PluginFuncs{
		Name:   name,
		Report: func(node ast.Node) string { return "issue" },
		Traverse: func() types.Traverser {
			return types.Traverser{
				key: func(node ast.Node, push func(ast.Node)) {
					push(node)
				},
			}
		},
		Fix: func(node ast.Node, opts map[string]any) {},
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

	Test(t, "runner: replacer returns 1 place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "runner: replacer place has correct rule", func(t *T) {
		t.Equal(places[0].Rule, "eq")
		t.End()
	})

	Test(t, "runner: replacer place has correct message", func(t *T) {
		t.Equal(places[0].Message, "use DeepEqual")
		t.End()
	})
}

func TestRunReplacerFixRewrites(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: replacer fix rewrites to DeepEqual", func(t *T) {
		t.Ok(strings.Contains(out, "t.DeepEqual(a, b)"))
		t.End()
	})
}

func TestRunMsgOverridesReport(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	kinds := loader.Load(replacerItem(), loader.Config{})
	pl := []PluginItem{{Rule: "eq", Plugin: kinds[0], Msg: "custom"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: msg overrides report", func(t *T) {
		t.Equal(places[0].Message, "custom")
		t.End()
	})
}

func TestRunTraverserReturnsPlaces(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	pl := items([]loader.PluginFuncs{traverserFuncs("file", "*ast.File")})
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: traverser returns 1 place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})
}

func TestRunTraverserFixCallsFix(t *testing.T) {
	src := "package p\n\nfunc f() {}\n"
	file, fset := parse(t, src)
	funcs := loader.PluginFuncs{
		Name:   "file",
		Report: func(node ast.Node) string { return "file issue" },
		Traverse: func() types.Traverser {
			return types.Traverser{
				"*ast.File": func(node ast.Node, push func(ast.Node)) {
					push(node)
				},
			}
		},
		Fix: func(node ast.Node, opts map[string]any) {
			f := node.(*ast.File)
			f.Name.Name = "q"
		},
	}
	pl := items([]loader.PluginFuncs{funcs})
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})

	Test(t, "runner: traverser fix renames package", func(t *T) {
		t.Equal(file.Name.Name, "q")
		t.End()
	})
}

func TestRunTraverserBlockVisitor(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	pl := items([]loader.PluginFuncs{traverserFuncs("block", "*ast.BlockStmt")})
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: block traverser returns 1 place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})
}

func TestRunEmptyPlugins(t *testing.T) {
	src := "package p\n\nfunc f() {}\n"
	file, fset := parse(t, src)
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: nil})

	Test(t, "runner: empty plugins return no places", func(t *T) {
		result := len(places)
		t.Equal(result, 0)

		t.End()
	})
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

	Test(t, "runner: guard rejecting match yields no places", func(t *T) {
		result := len(places)
		t.Equal(result, 0)

		t.End()
	})
}

func TestRunNilGuardPanics(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "nilguard",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{"t.Equal(__a, __b)": nil} },
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "x"} },
	}}
	pl := items(funcs)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil guard")
		}
	}()
	RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})
}

func TestRunTraverserBlockFix(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	funcs := loader.PluginFuncs{
		Name:   "block",
		Report: func(node ast.Node) string { return "block" },
		Traverse: func() types.Traverser {
			return types.Traverser{
				"*ast.BlockStmt": func(node ast.Node, push func(ast.Node)) {
					push(node)
				},
			}
		},
		Fix: func(node ast.Node, opts map[string]any) {
			block := node.(*ast.BlockStmt)
			block.List = nil
		},
	}
	pl := items([]loader.PluginFuncs{funcs})
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	fileBlock := file.Decls[0].(*ast.FuncDecl).Body

	Test(t, "runner: block fix empties stmt list", func(t *T) {
		result := len(fileBlock.List)
		t.Equal(result, 0)

		t.End()
	})
}

func TestRunTraverserMsgOverride(t *testing.T) {
	src := "package p\n\nfunc f() {}\n"
	file, fset := parse(t, src)
	kinds := loader.Load([]loader.PluginFuncs{traverserFuncs("file", "*ast.File")}, loader.Config{})
	pl := []PluginItem{{Rule: "file", Plugin: kinds[0], Msg: "override"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: traverser msg override", func(t *T) {
		t.Equal(places[0].Message, "override")
		t.End()
	})
}

func TestRunFixCountLoops(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n\tt.Equal(c, d)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 3, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: fix count loops converts both", func(t *T) {
		result := strings.Count(out, "DeepEqual")
		t.Equal(result, 2)

		t.End()
	})
}

func TestRunNoFixNoRewrite(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	pl := items(replacerItem())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: false, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: no fix means no rewrite", func(t *T) {
		t.NotOk(strings.Contains(out, "DeepEqual"))
		t.End()
	})
}

func TestSubstituteUnmatchedHole(t *testing.T) {
	got := substitute("hello __nope", Vars{})

	Test(t, "runner: unmatched hole preserved", func(t *T) {
		t.Equal(got, "hello __nope")
		t.End()
	})
}

func TestSubstituteMatchedHole(t *testing.T) {
	src := "package p\n\nfunc f() { x := 1 }\n"
	file, _ := parse(t, src)
	stmt := file.Decls[0].(*ast.FuncDecl).Body.List[0]
	got := substitute("a = __x", Vars{"__x": stmt})

	Test(t, "runner: matched hole substituted", func(t *T) {
		t.Ok(strings.Contains(got, "x := 1"))
		t.End()
	})
}

func TestRunRenderArgSlice(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tuse(f, g)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "args",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{"use(__args)": func(v types.Vars) bool { return true }} },
		Replace: func() types.Replacer { return types.Replacer{"use(__args)": "wrap(__args)"} },
	}}
	pl := items(funcs)
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: renders wrapped arg slice", func(t *T) {
		t.Ok(strings.Contains(out, "wrap(f, g)"))
		t.End()
	})
}

func TestRunRenderBodySlice(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tTest(\"s\", \"n\", func(t *T) {\n\t\tt.Equal(a, b)\n\t})\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:   "body",
		Report: func() string { return "m" },
		Match: func() types.Matcher {
			return types.Matcher{"Test(__a, __b, func(__r *T) { __body })": func(v types.Vars) bool { return true }}
		},
		Replace: func() types.Replacer {
			return types.Replacer{"Test(__a, __b, func(__r *T) { __body })": "Test(__a, __b, func(__r *T) {\n__body\n__r.End()\n})"}
		},
	}}
	pl := items(funcs)
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: renders body with End added", func(t *T) {
		t.Ok(strings.Contains(out, ".End()"))
		t.End()
	})
}

func TestPrintNodeNil(t *testing.T) {
	got := printNode(nil)

	Test(t, "runner: nil node prints empty", func(t *T) {
		t.Equal(got, "")
		t.End()
	})
}

func TestApplyRewritesMultipleInBlock(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tmakeSlices(x)\n\tmakeSlices(y)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "multi",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{"makeSlices(__x)": func(v types.Vars) bool { return true }} },
		Replace: func() types.Replacer { return types.Replacer{"makeSlices(__x)": "v := __x"} },
	}}
	pl := items(funcs)
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: two rewrites in block", func(t *T) {
		result := strings.Count(out, "v := ")
		t.Equal(result, 2)

		t.End()
	})
}

func TestRunUnparseableReplace(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:   "badreplace",
		Report: func() string { return "m" },
		Match: func() types.Matcher {
			return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}
		},
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "func ("} },
	}}
	pl := items(funcs)
	places := RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})

	Test(t, "runner: unparseable replace still reports place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})
}

func TestRunTraverserBlockMsgOverride(t *testing.T) {
	src := "package p\n\nfunc f() {\n\tx := 1\n}\n"
	file, fset := parse(t, src)
	kinds := loader.Load([]loader.PluginFuncs{traverserFuncs("block", "*ast.BlockStmt")}, loader.Config{})
	pl := []PluginItem{{Rule: "block", Plugin: kinds[0], Msg: "blockmsg"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: block traverser msg override", func(t *T) {
		t.Equal(places[0].Message, "blockmsg")
		t.End()
	})
}

func TestSubstituteAndParseError(t *testing.T) {
	stmts := substituteAndParse("func (", Vars{})

	Test(t, "runner: unparseable template returns nil", func(t *T) {
		t.NotOk(stmts)

		t.End()
	})
}

// declFuncs builds a replacer plugin whose Match and Replace carry a
// top-level declaration pattern (used for decl-level rewrites).
func declFuncs() []loader.PluginFuncs {
	return []loader.PluginFuncs{{
		Name:   "decl",
		Report: func() string { return "decl issue" },
		Match: func() types.Matcher {
			return types.Matcher{`func Match() Matcher { return Matcher{__a: nil} }`: func(v types.Vars) bool { return true }}
		},
		Replace: func() types.Replacer { return types.Replacer{`func Match() Matcher { return Matcher{__a: nil} }`: ""} },
	}}
}

func TestRunDeclRewritesReportsPlace(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	pl := items(declFuncs())
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: decl rewrite reports 1 place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})
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

	Test(t, "runner: decl guard rejecting yields no places", func(t *T) {
		result := len(places)
		t.Equal(result, 0)

		t.End()
	})
}

func TestRunDeclRewritesRemovesDecl(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	pl := items(declFuncs())
	RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: decl rewrite removes declaration", func(t *T) {
		t.NotOk(strings.Contains(out, "func Match"))
		t.End()
	})
}

func TestRunDeclRewritesNoMatch(t *testing.T) {
	src := "package p\n\nfunc Other() {}\n"
	file, fset := parse(t, src)
	pl := items(declFuncs())
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: decl non-match yields no places", func(t *T) {
		result := len(places)
		t.Equal(result, 0)

		t.End()
	})
}

func TestRunDeclRewritesMsgOverride(t *testing.T) {
	src := "package p\n\nfunc Match() Matcher { return Matcher{\"x\": nil} }\n"
	file, fset := parse(t, src)
	kinds := loader.Load(declFuncs(), loader.Config{})
	pl := []PluginItem{{Rule: "decl", Plugin: kinds[0], Msg: "custom decl"}}
	places := RunPlugins(RunParams{File: file, Fset: fset, Plugins: pl})

	Test(t, "runner: decl msg override", func(t *T) {
		t.Equal(places[0].Message, "custom decl")
		t.End()
	})
}

func TestApplyDeclRewritesKeepsNonEmptyTmpl(t *testing.T) {
	file, _ := parse(t, "package p\n\nfunc Match() Matcher { return Matcher{} }\n")
	applyDeclRewrites(file, []declRewrite{{idx: 0, tmpl: "func Other() {}"}})

	Test(t, "runner: non-empty decl template kept", func(t *T) {
		result := len(file.Decls)
		t.Equal(result, 1)

		t.End()
	})
}

func TestApplyDeclRewritesRemovesMultiple(t *testing.T) {
	file, _ := parse(t, "package p\n\nfunc A() {}\n\nfunc B() {}\n")
	applyDeclRewrites(file, []declRewrite{{idx: 0, tmpl: ""}, {idx: 1, tmpl: ""}})

	Test(t, "runner: multiple empty templates removed", func(t *T) {
		result := len(file.Decls)
		t.Equal(result, 0)

		t.End()
	})
}

func TestRunMatchOnlyFixNoRewrite(t *testing.T) {
	src := "package p\n\nfunc f() { t.Equal(a, b) }\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:   "report-only",
		Report: func() string { return "m" },
		Match: func() types.Matcher {
			return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars) bool { return true }}
		},
		Replace: func() types.Replacer { return types.Replacer{} },
	}}
	pl := items(funcs)
	places := RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: match-only reports place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "runner: match-only does not rewrite", func(t *T) {
		t.NotOk(strings.Contains(out, "DeepEqual"))
		t.End()
	})
}

func TestRunReplaceOnlyRewrites(t *testing.T) {
	src := "package p\n\nfunc f() { t.Equal(a, b) }\n"
	file, fset := parse(t, src)
	funcs := []loader.PluginFuncs{{
		Name:    "replace-only",
		Report:  func() string { return "m" },
		Match:   func() types.Matcher { return types.Matcher{} },
		Replace: func() types.Replacer { return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"} },
	}}
	pl := items(funcs)
	places := RunPlugins(RunParams{File: file, Fset: fset, Fix: true, FixCount: 1, Plugins: pl})
	out := printFile(t, file, fset)

	Test(t, "runner: replace-only reports place", func(t *T) {
		result := len(places)
		t.Equal(result, 1)

		t.End()
	})

	Test(t, "runner: replace-only rewrites", func(t *T) {
		t.Ok(strings.Contains(out, "DeepEqual"))
		t.End()
	})
}

func TestRunInjectsBlockGuard(t *testing.T) {
	src := "package p\nfunc f() { t.Equal(1, 2) }\n"
	file, fset := parse(t, src)
	guardCalled := false
	sawBlock := false
	pf := loader.PluginFuncs{
		Name:   "block-guard",
		Report: func() string { return "found" },
		Match: func() types.Matcher {
			return types.Matcher{
				"t.Equal(__a, __b)": func(vars types.Vars) bool {
					guardCalled = true
					_, sawBlock = vars["$block"].(*ast.BlockStmt)
					return true
				},
			}
		},
		Replace: func() types.Replacer {
			return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}
		},
	}
	RunPlugins(RunParams{File: file, Fset: fset, Plugins: items([]loader.PluginFuncs{pf})})

	Test(t, "runner: $block injected into vars before guard", func(t *T) {
		t.Ok(guardCalled && sawBlock)
		t.End()
	})
}
