# 🐘 Indra [![Build Status](https://github.com/coderaiser/indra/actions/workflows/go.yml/badge.svg)](https://github.com/coderaiser/indra/actions) [![Coverage Status](https://coveralls.io/repos/github/coderaiser/indra/badge.svg?branch=master)](https://coveralls.io/github/coderaiser/indra?branch=master)

🐘 **Indra** is a Go port of 🐊[**Putout**](https://github.com/coderaiser/putout) — a pluggable, configurable linter and code transformer. Write declarative codemods in Go in the simplest possible way.

## Table of contents

- [🚚 Installation](#-installation)
- [🎙 Usage](#-usage)
- [🏛 Architecture](#-architecture)
- [🏗 API](#-api)
- [⚙️ Configuration](#️-configuration)
- [🏟 Plugins](#-plugins)

## 🚚 Installation

```sh
go install coderaiser/indra/cmd/indra@latest
```

Or with [palabra](https://github.com/coderaiser/palabra):

```sh
palabra i indra
```

## 🎙 Usage

Find issues in a directory:

```sh
indra ./...
```

Apply fixes automatically:

```sh
indra --fix ./...
```

## 🏛 Architecture

🐘 **Indra** follows the same architecture as 🐊 **Putout**:

```
Source file
    │
    ▼
engine-parser   — parse Go source into *ast.File
    │
    ▼
engine-runner   — run plugins (match → replace, traverse)
    │
    ▼
engine-printer  — print *ast.File back to source
    │
    ▼
engine-processor — coordinate parser/runner/printer per file type
    │
    ▼
engine-loader   — resolve .indra.toml chain → Options
```

### 🌴 Laws of the Jungle

- 🐅 engines chase plugins and processors;
- 🦌 plugins know nothing about each other;
- 🦒 processors handle file types;

### 💚 Engines

| Package | Description |
| --- | --- |
| `engine-parser` | `Parse(src []byte) (*ast.File, *token.FileSet, error)` |
| `engine-runner` | `Run(file, fset, plugins, fix) ([]Place, error)` |
| `engine-printer` | `Print(file *ast.File, fset *token.FileSet) ([]byte, error)` |
| `engine-processor` | `Process(name string, src []byte, opts Options) ([]byte, []Place, error)` |
| `engine-loader` | `Load(dir string, defaults []byte) (Options, error)` |

### 🧪 Processors

| Package | Handles |
| --- | --- |
| `processor-go` | `*.go` files |

## 🏗 API

```go
import indra "coderaiser/indra"

// Lint or fix source. nil Fix means fix is enabled.
out, places, err := indra.Indra(src, indra.Options{
    Processors: ...,
    PluginList: ...,
    Fix:        nil, // nil = fix enabled
})

// Load config from .indra.toml chain walking up from dir.
opts, err := indra.Load(".")
```

#### Report only (no fix)

```go
fix := false
out, places, err := indra.Indra(src, indra.Options{
    PluginList: plugins,
    Fix:        &fix,
})
// places contains all findings; out == src
```

#### Fix

```go
out, places, err := indra.Indra(src, indra.Options{
    PluginList: plugins,
    Fix:        nil, // nil = enabled
})
// out contains fixed source
```

## ⚙️ Configuration

Indra reads `.indra.toml` files walking up from the linted file's directory and uses the same merge strategy as 🐊 **Putout** uses with `.putout.json`.

**`cmd/indra/indra.toml`** — CLI defaults (embedded, all rules on):

```toml
processors = ["go"]
plugins    = ["tape", "remove-unused-import", "remove-unused-variable"]

[rules]
"remove-unused-variable"             = "off"
```

**`.indra.toml`** — repo-level overrides (highest priority for user config):

```toml
[rules]
"tape/remove-skip" = "off"
```

Merge order (highest priority first):

1. deepest `.indra.toml` in the directory tree
2. parent `.indra.toml` files walking up to repo root
3. `cmd/indra/indra.toml` CLI defaults

## 🍄 License

MIT
