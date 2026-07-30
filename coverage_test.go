package indra_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	indra "coderaiser/indra"

	tape "github.com/coderaiser/go-tape"

	"github.com/lithammer/dedent"
)

func writeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprint(f, content)
	return err
}

func TestParseCoverage(t *testing.T) {
	tape.Test(t, "indra: parse returns uncovered blocks", func(t *tape.T) {
		input := dedent.Dedent(`
            mode: set
            github.com/app/main.go:5.1,8.2 3 1
            github.com/app/main.go:10.1,12.2 2 0
        `)
		blocks := indra.ParseCoverage(strings.NewReader(input))
		t.DeepEqual(blocks, []indra.Block{
			{File: "github.com/app/main.go", Start: 10, End: 12},
		})
		t.End()
	})

	tape.Test(t, "indra: parse skips covered blocks", func(t *tape.T) {
		input := dedent.Dedent(`
            mode: set
            github.com/app/main.go:1.1,2.1 1 5
        `)
		blocks := indra.ParseCoverage(strings.NewReader(input))
		t.DeepEqual(blocks, []indra.Block(nil))
		t.End()
	})

	tape.Test(t, "indra: parse returns nil on empty input", func(t *tape.T) {
		blocks := indra.ParseCoverage(strings.NewReader("mode: set\n"))
		t.DeepEqual(blocks, []indra.Block(nil))
		t.End()
	})
}

func TestFormatBlock(t *testing.T) {
	tape.Test(t, "indra: format block without lines", func(t *tape.T) {
		result := indra.FormatBlock(
			indra.Block{File: "main.go", Start: 10, End: 12},
			"/", nil, false,
		)
		t.Equal(result, "file://main.go:10: 10-12")
		t.End()
	})

	tape.Test(t, "indra: format block with lines contains line prefix", func(t *tape.T) {
		lines := []string{"if x == nil {", "    return err", "}"}
		result := indra.FormatBlock(
			indra.Block{File: "main.go", Start: 10, End: 12},
			"/", lines, false,
		)
		t.Match(result, "10 | if x == nil {")
		t.End()
	})

	tape.Test(t, "indra: format block with color contains ANSI code", func(t *tape.T) {
		lines := []string{"return nil"}
		result := indra.FormatBlock(
			indra.Block{File: "main.go", Start: 5, End: 5},
			"/", lines, true,
		)
		t.Match(result, "\033[31m")
		t.End()
	})

	tape.Test(t, "indra: format block line number 20", func(t *tape.T) {
		lines := []string{"a", "b", "c"}
		result := indra.FormatBlock(
			indra.Block{File: "f.go", Start: 20, End: 22},
			"/", lines, false,
		)
		t.Match(result, "20")
		t.End()
	})

	tape.Test(t, "indra: format block line number 21", func(t *tape.T) {
		lines := []string{"a", "b", "c"}
		result := indra.FormatBlock(
			indra.Block{File: "f.go", Start: 20, End: 22},
			"/", lines, false,
		)
		t.Match(result, "21")
		t.End()
	})

	tape.Test(t, "indra: format block line number 22", func(t *tape.T) {
		lines := []string{"a", "b", "c"}
		result := indra.FormatBlock(
			indra.Block{File: "f.go", Start: 20, End: 22},
			"/", lines, false,
		)
		t.Match(result, "22")
		t.End()
	})
}

func TestReadLines(t *testing.T) {
	tape.Test(t, "indra: ReadLines returns correct range", func(t *tape.T) {
		path := t.TB().TempDir() + "/test.go"
		if err := writeFile(path, "line1\nline2\nline3\nline4\nline5\n"); err != nil {
			t.TB().Fatal(err)
		}
		lines, _ := indra.ReadLines(path, 2, 4)
		t.DeepEqual(lines, []string{"line2", "line3", "line4"})
		t.End()
	})

	tape.Test(t, "indra: ReadLines returns error on missing file", func(t *tape.T) {
		_, err := indra.ReadLines("/nonexistent/file.go", 1, 5)
		t.Error(err)
		t.End()
	})
}

func TestColorEnabled(t *testing.T) {
	tape.Test(t, "indra: ColorEnabled returns true when COLOR=1", func(t *tape.T) {
		t.Setenv("COLOR", "1")
		t.Ok(indra.ColorEnabled())
		t.End()
	})

	tape.Test(t, "indra: ColorEnabled returns false when COLOR=0", func(t *tape.T) {
		t.Setenv("COLOR", "0")
		t.NotOk(indra.ColorEnabled())
		t.End()
	})
}

func TestHighlightLines(t *testing.T) {
	tape.Test(t, "indra: HighlightLines returns ANSI codes", func(t *tape.T) {
		lines := []string{"func main() {", "\treturn", "}"}
		result := indra.HighlightLines(lines)
		t.Match(strings.Join(result, "\n"), "\033[")
		t.End()
	})

	tape.Test(t, "indra: HighlightLines preserves line count", func(t *tape.T) {
		lines := []string{"func main() {", "\treturn", "}"}
		result := indra.HighlightLines(lines)
		t.Equal(len(result), len(lines))
		t.End()
	})

	tape.Test(t, "indra: HighlightLines returns fallback on empty input", func(t *tape.T) {
		result := indra.HighlightLines([]string{})
		t.DeepEqual(result, []string{""})
		t.End()
	})
}

func TestFindModule(t *testing.T) {
	tape.Test(t, "indra: FindModule returns root dir", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(filepath.Join(dir, "go.mod"), "module mymod/myapp\n\ngo 1.22\n")
		root, _ := indra.FindModule(dir)
		t.Equal(root, dir)
		t.End()
	})

	tape.Test(t, "indra: FindModule returns module name", func(t *tape.T) {
		dir := t.TB().TempDir()
		writeFile(filepath.Join(dir, "go.mod"), "module mymod/myapp\n\ngo 1.22\n")
		_, name := indra.FindModule(dir)
		t.Equal(name, "mymod/myapp")
		t.End()
	})
}

func TestRelativeFile(t *testing.T) {
	tape.Test(t, "indra: RelativeFile strips module prefix", func(t *tape.T) {
		result := indra.RelativeFile("mymod/myapp/pkg/foo.go", "mymod/myapp")
		t.Equal(result, "pkg/foo.go")
		t.End()
	})

	tape.Test(t, "indra: RelativeFile returns path unchanged when no match", func(t *tape.T) {
		result := indra.RelativeFile("other/module/foo.go", "mymod/myapp")
		t.Equal(result, "other/module/foo.go")
		t.End()
	})
}

func TestResolveFile(t *testing.T) {
	tape.Test(t, "indra: ResolveFile strips module name from path", func(t *tape.T) {
		dir := t.TB().TempDir()
		if err := writeFile(filepath.Join(dir, "go.mod"), "module mymod/myapp\n\ngo 1.22\n"); err != nil {
			t.TB().Fatal(err)
		}
		result := indra.ResolveFile("pkg/foo.go", dir)
		t.Equal(result, filepath.Join(dir, "pkg/foo.go"))
		t.End()
	})

	tape.Test(t, "indra: ResolveFile returns path unchanged when no module", func(t *tape.T) {
		result := indra.ResolveFile("some/path/foo.go", t.TB().TempDir())
		t.Equal(result, "some/path/foo.go")
		t.End()
	})
}

func TestMergeBlocks(t *testing.T) {
	tape.Test(t, "indra: MergeBlocks merges overlapping same-file blocks", func(t *tape.T) {
		result := indra.MergeBlocks([]indra.Block{
			{File: "a.go", Start: 10, End: 10},
			{File: "a.go", Start: 10, End: 12},
			{File: "a.go", Start: 13, End: 15},
		})
		t.DeepEqual(result, []indra.Block{
			{File: "a.go", Start: 10, End: 15},
		})
		t.End()
	})

	tape.Test(t, "indra: MergeBlocks keeps different files separate", func(t *tape.T) {
		result := indra.MergeBlocks([]indra.Block{
			{File: "b.go", Start: 1, End: 1},
			{File: "a.go", Start: 1, End: 1},
		})
		t.DeepEqual(result, []indra.Block{
			{File: "a.go", Start: 1, End: 1},
			{File: "b.go", Start: 1, End: 1},
		})
		t.End()
	})
}
