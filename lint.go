package indra

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"

	loader "coderaiser/indra/engine_loader"
	processor "coderaiser/indra/engine_processor"
	runner "coderaiser/indra/engine_runner"
	"coderaiser/indra/internal/config"
	"coderaiser/indra/internal/formatter"
	"coderaiser/indra/internal/formatter_progress_bar"
	processor_go "coderaiser/indra/processor-go"
	"coderaiser/indra/types"
)

// ErrInvalidOption is returned when an unknown flag is passed to Indra.
var ErrInvalidOption = errors.New("invalid option")

//go:embed cmd/indra/help.toml
var helpToml []byte

type helpConfig struct {
	Usage   struct{ Text string }
	Options []struct {
		Flag string
		Desc string
		Note string
	}
}

func loadUsage() string { return parseUsage(helpToml) }

func parseUsage(data []byte) string {
	var cfg helpConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return "Usage: indra [options] [path ...]\n"
	}
	var sb strings.Builder
	sb.WriteString(cfg.Usage.Text + "\n\nOptions:\n")
	for _, o := range cfg.Options {
		fmt.Fprintf(&sb, "  %-24s %s\n", o.Flag, o.Desc)
		if o.Note != "" {
			fmt.Fprintf(&sb, "                         %s\n", o.Note)
		}
	}
	sb.WriteString("\n")
	return sb.String()
}

// loadPlugins resolves all enabled plugins into runnable items.
func loadPlugins(registry []loader.PluginFuncs, lc loader.Config) []runner.PluginItem {
	kinds := loader.Load(registry, lc)
	items := make([]runner.PluginItem, 0, len(kinds))
	for _, k := range kinds {
		items = append(items, runner.PluginItem{Rule: k.Name(), Plugin: k})
	}
	return items
}

// filterPlugins restricts items to those named in the [plugins] list. A list
// entry matches a rule exactly or, for a group name like "tape", all of its
// "tape/*" sub-rules. items is filtered in place.
func filterPlugins(items []runner.PluginItem, names []string) []runner.PluginItem {
	out := items[:0]
	for _, item := range items {
		for _, name := range names {
			if item.Rule == name || strings.HasPrefix(item.Rule, name+"/") {
				out = append(out, item)
				break
			}
		}
	}
	return out
}

// configForFile returns the effective loader.Config for a filename: the global
// [rules] merged with any [match] overrides that apply to that file.
func configForFile(cfg config.Config, filename string) loader.Config {
	lc := cfg.ToLoaderConfig()
	for rule, val := range cfg.Match.OverrideRules(filename) {
		lc[rule] = loader.RuleState{Enabled: val == "on"}
	}
	return lc
}

// Lint runs all plugins against src.
// Returns rewritten source, findings, and any parse error.
func Lint(registry []loader.PluginFuncs, src []byte, fix bool) ([]byte, []types.Place, error) {
	res, err := processor.Process(processor.Params{
		Src:     src,
		Fix:     fix,
		Plugins: loadPlugins(registry, loader.DefaultConfig()),
	})
	return res.Out, res.Places, err
}

// Indra runs the CLI over the given files, printing findings.
func Indra(registry []loader.PluginFuncs, args []string, w io.Writer) error {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			fmt.Fprint(w, loadUsage())
			return nil
		}
		if arg == "--version" || arg == "-v" {
			fmt.Fprintln(w, VersionLine())
			return nil
		}
	}

	fix := false
	formatterName := os.Getenv("INDRA_FORMATTER")
	files := []string{}
	var unknownFlags []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--fix":
			fix = true
		case a == "-f", a == "--format":
			if i+1 < len(args) {
				i++
				formatterName = args[i]
			}
		case strings.HasPrefix(a, "-f="):
			formatterName = strings.TrimPrefix(a, "-f=")
		case strings.HasPrefix(a, "--format="):
			formatterName = strings.TrimPrefix(a, "--format=")
		case strings.HasPrefix(a, "-"):
			unknownFlags = append(unknownFlags, a)
		default:
			files = append(files, a)
		}
	}

	if len(unknownFlags) > 0 {
		fmt.Fprintf(w, "🐊 Invalid option `%s`.\n", unknownFlags[0])
		return fmt.Errorf("%w: %s", ErrInvalidOption, unknownFlags[0])
	}

	if len(files) == 0 {
		return nil
	}

	cfg, _ := config.Load(".")
	ignore := append(slices.Clone(config.DefaultIgnorePatterns), cfg.Ignore.Patterns...)
	formatter_progress_bar.Configure(formatter_progress_bar.Config{
		Color:    cfg.Progress.Color,
		MinCount: cfg.Progress.MinCount,
	})

	rawFiles, dirs := processor_go.ResolveArgs(files)
	allFiles := processor_go.CollectFiles(rawFiles, dirs, ignore)
	if len(allFiles) == 0 {
		return nil
	}

	format := formatter.ChooseByName(formatterName)
	failed := false
	filesWithIssues := 0
	errorsCount := 0
	total := len(allFiles)

	for i, filename := range allFiles {
		fileItems := loadPlugins(registry, configForFile(cfg, filename))
		if len(cfg.Plugins) > 0 {
			fileItems = filterPlugins(fileItems, cfg.Plugins)
		}
		src, readErr := os.ReadFile(filename)
		places, err := processor_go.ProcessFile(filename, processor_go.Opt(fileItems, fix))
		if err != nil {
			fmt.Fprintf(w, "file://%s: %v\n", filename, err)
			failed = true
			continue
		}
		_ = readErr
		if len(places) > 0 {
			filesWithIssues++
			errorsCount += len(places)
			failed = true
		}
		out := format(filename, src, places, i, total, filesWithIssues, errorsCount)
		if out != "" {
			fmt.Fprint(w, out)
		}
	}

	if failed {
		return errors.New("lint failed")
	}
	return nil
}
