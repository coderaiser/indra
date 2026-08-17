package engine_processor

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"
	"testing"

	loader "coderaiser/indra/engine_loader"
	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/types"
)

// eqReplacer is a synthetic replacer matching Equal calls.
type eqReplacer struct{}

func (eqReplacer) Report() string { return "use DeepEqual" }
func (eqReplacer) Match() types.Matcher {
	return types.Matcher{"t.Equal(__a, __b)": func(v types.Vars, _ types.Path) bool { return true }}
}
func (eqReplacer) Replace() types.Replacer {
	return types.Replacer{"t.Equal(__a, __b)": "t.DeepEqual(__a, __b)"}
}

func pluginItems() []runner.PluginItem {
	kinds := loader.Load([]loader.PluginFuncs{{Name: "eq", Plugin: eqReplacer{}}}, loader.Config{})
	return []runner.PluginItem{{Rule: kinds[0].Name(), Plugin: kinds[0]}}
}

func TestProcessReturnsPlacesWithoutFix(t *testing.T) {
	src := []byte("package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	res, err := Process(Params{Src: src, Plugins: pluginItems()})
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	if len(res.Places) != 1 {
		t.Fatalf("expected 1 place, got %d", len(res.Places))
	}
	if res.Places[0].Rule != "eq" {
		t.Fatalf("unexpected rule %q", res.Places[0].Rule)
	}
	if string(res.Out) != string(src) {
		t.Fatal("without fix, out must equal src")
	}
}

func TestProcessFixes(t *testing.T) {
	src := []byte("package p\n\nfunc f() {\n\tt.Equal(a, b)\n}\n")
	res, err := Process(Params{Src: src, Fix: true, Plugins: pluginItems()})
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	if !strings.Contains(string(res.Out), "t.DeepEqual(a, b)") {
		t.Fatalf("expected rewritten output:\n%s", res.Out)
	}
}

func TestProcessParseError(t *testing.T) {
	src := []byte("package p\nfunc (\n")
	_, err := Process(Params{Src: src, Plugins: pluginItems()})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestProcessNoPlugins(t *testing.T) {
	src := []byte("package p\n\nfunc f() {}\n")
	res, err := Process(Params{Src: src, Plugins: nil})
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	if len(res.Places) != 0 {
		t.Fatalf("expected no places, got %d", len(res.Places))
	}
	if string(res.Out) != string(src) {
		t.Fatal("expected src unchanged")
	}
}

func TestProcessFixNoChange(t *testing.T) {
	src := []byte("package p\n\nfunc f() {}\n")
	res, err := Process(Params{Src: src, Fix: true, Plugins: pluginItems()})
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	if string(res.Out) != string(src) {
		t.Fatalf("expected unchanged output for clean file:\n%s", res.Out)
	}
}

func TestProcessPrinterError(t *testing.T) {
	orig := print
	print = func(file *ast.File, fset *token.FileSet) ([]byte, error) {
		return nil, errors.New("print failed")
	}
	defer func() { print = orig }()

	src := []byte("package p\n\nfunc f() {}\n")
	res, err := Process(Params{Src: src, Fix: true, Plugins: pluginItems()})
	if err == nil {
		t.Fatal("expected printer error")
	}
	if string(res.Out) != string(src) {
		t.Fatal("expected src preserved on printer error")
	}
}
