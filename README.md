# 🐘 Indra

🐘 **Indra** is a Go port of 🐊[**Putout**](https://github.com/coderaiser/putout) —
a pluggable, configurable linter and code transformer for Go test files
using [📼 go-tape](https://github.com/coderaiser/go-tape).

## 🚚 Installation

    go install coderaiser/indra/cmd/indra@latest

## 🎙 Usage

    indra ./...        # report issues
    indra --fix ./...  # fix issues

## 🏛 Architecture

    Source
        │
        ▼
    engine_parser      — parse Go source into *ast.File
        │
        ▼
    engine_runner      — run plugins (Match → Replace, Traverse)
        │
        ▼
    engine_printer     — print *ast.File back to []byte
        │
        ▼
    engine_processor   — orchestrate parser/runner/printer per file type
        │
        ▼
    engine_loader      — resolve .indra.toml chain → Options

## ⚙️ Configuration

Create `.indra.toml` in your project root:

    processors = ["go"]
    plugins    = ["tape", "remove-unused-import", "remove-unused-variable"]

    [rules]
    "tape/remove-skip" = "off"

    [match]
    "*_integration_test.go" = { "tape/add-t-end" = "off" }

## 🏟 Plugins

### 📼 tape

| Rule | Description |
|---|---|
| `tape/remove-skip` | Replace `Test.Skip` with `Test` |
| `tape/add-t-end` | Add missing `t.End()` |
| `tape/convert-equal-to-deep-equal` | Use `DeepEqual` for slice args |
| `tape/extract-result-from-assertion` | Extract inline expressions from assertions |

### 🔧 Single rules

| Rule | Description |
|---|---|
| `remove-unused-import` | Remove unused imports |
| `remove-unused-variable` | Remove unused variables |

## 🍄 License

MIT
