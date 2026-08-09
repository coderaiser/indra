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

Rules for [go-tape](https://github.com/coderaiser/go-tape) test files.
Enabled for `*_test.go` by default.

| Rule | README | Description |
|---|---|---|
| `tape/remove-skip` | [📖](internal/plugin_tape/remove_skip/README.md) | Replace `Test.Skip` with `Test` |
| `tape/add-t-end` | [📖](internal/plugin_tape/add_t_end/README.md) | Add missing `t.End()` |
| `tape/apply-dedent` | [📖](internal/plugin_tape/apply_dedent/README.md) | Remove `dedent.Dedent` wrappers |
| `tape/convert-equal-to-deep-equal` | [📖](internal/plugin_tape/convert_equal_to_deep_equal/README.md) | Use `DeepEqual` for slice args |
| `tape/convert-equal-to-not-ok` | [📖](internal/plugin_tape/convert_equal_to_not_ok/README.md) | Convert `Equal(x, nil/false)` to `NotOk(x)` |
| `tape/convert-equal-to-ok` | [📖](internal/plugin_tape/convert_equal_to_ok/README.md) | Convert `Equal(x, true)` to `Ok(x)` |
| `tape/convert-no-error-to-not-ok` | [📖](internal/plugin_tape/convert_no_error_to_not_ok/README.md) | Convert `NoError` to `NotOk` |
| `tape/convert-ok-to-not-ok` | [📖](internal/plugin_tape/convert_ok_to_not_ok/README.md) | Convert `Ok(err == nil)` to `NotOk(err)` |
| `tape/extract-result-from-assertion` | [📖](internal/plugin_tape/extract_result_from_assertion/README.md) | Extract inline call from assertion |
| `tape/remove-useless-condition` | [📖](internal/plugin_tape/remove_useless_condition/README.md) | Remove useless condition in `Ok`/`NotOk` |
| `tape/remove-useless-prefix` | [📖](internal/plugin_tape/remove_useless_prefix/README.md) | Remove redundant tape qualifier |

### 🚦 conditions

General Go code-quality rules. Enabled for all files by default.

| Rule | README | Description |
|---|---|---|
| `conditions/remove-useless-comments` | [📖](internal/plugin_conditions/remove_useless_comments/README.md) | Remove separator banner comments |
| `conditions/convert-switch-to-if` | [📖](internal/plugin_conditions/convert_switch_to_if/README.md) | Replace qualifying `switch` with `if` chains |

### 🔧 indra

Meta-rules for indra plugin authoring. Enabled via `"indra" = "on"`.

| Rule | README | Default | Description |
|---|---|---|---|
| `indra/remove-useless-match` | [📖](internal/plugin_indra/remove_useless_match/README.md) | off | Remove nil/empty `Match()` entries |
| `indra/apply-compare` | [📖](internal/plugin_indra/apply_compare/README.md) | on | Use `Compare` over `GetTemplateValues != nil` |
| `indra/apply-fixture-name-to-message` | [📖](internal/plugin_indra/apply_fixture_name_to_message/README.md) | on | Prefix test message with rule name |
| `indra/replace-test-message` | [📖](internal/plugin_indra/replace_test_message/README.md) | on | Include fixture name in test message |
| `indra/convert-for-to-create-test` | [📖](internal/plugin_indra/convert_for_to_create/README.md) | off | Rename `For` to `CreateTest` |
| `indra/convert-inspect-to-traverse` | [📖](internal/plugin_indra/convert_inspect_to_traverse/README.md) | on | Replace `ast.Inspect` with `path.Traverse` |

### 📦 Single rules

| Rule | README | Description |
|---|---|---|
| `remove-unused-variables` | [📖](internal/plugin_remove_unused_variables/README.md) | Remove unused imports, variables, constants, private functions |

## 🍄 License

MIT
