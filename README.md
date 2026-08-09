# 🐘 Indra

<img width="500" height="320" alt="image" src="https://github.com/user-attachments/assets/b3fb1227-99c1-46bf-84d0-bb59c9c9f858" />


🐘**Indra** is a Go port of 🐊[**Putout**](https://github.com/coderaiser/putout) —
a pluggable, configurable linter and code transformer for Go.

## 🚚 Install

Install 🐘**Indra** with [palabra](https://github.com/coderaiser/palabra):

```
palabra i indra
```

Or go:

```sh
go install coderaiser/indra/cmd/indra@latest
```

## 🎙 Usage

```sh
indra .        # report issues
indra --fix .  # fix issues
```

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
| `tape/apply-dedent` | Remove `dedent.Dedent` wrappers |
| `tape/convert-equal-to-deep-equal` | Use `DeepEqual` for slice args |
| `tape/extract-result-from-assertion` | Extract inline expressions from assertions |
| `tape/remove-useless-condition` | Remove useless condition inside `Ok`/`NotOk` |

### 🚦 conditions

Rules enabled for all files — nothing tape-specific.

| Rule | Description |
|---|---|

### 🔧 Single rules

| Rule | Description |
|---|---|
| `remove-unused-import` | Remove unused imports |
| `remove-unused-variable` | Remove unused variables |

## 🍄 License

MIT
