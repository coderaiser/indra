# conditions

Rules that apply to any assertion library. Nothing here is tape-specific, so
the group is enabled for all files by default — no glob restriction needed.

## Rules

- ✅ [remove-useless-comments](#remove-useless-comments)
- ✅ [convert-switch-to-if](#convert-switch-to-if)

## remove-useless-comments

Removes separator banner comments built from repeated `─` characters.

### ❌ Incorrect

```go
package fixture

// ── imports ──────────────────────────────────────────────────────────────────

import "fmt"
```

### ✅ Correct

```go
package fixture

import "fmt"
```

## convert-switch-to-if

Use `if` instead of `switch` for one-value-per-case switches whose cases
end with a `return`.

### ❌ Incorrect

```go
func f(x string) string {
    switch x {
    case "a":
        return "A"
    case "b":
        return "B"
    }
    return "unknown"
}
```

### ✅ Correct

```go
func f(x string) string {
    if x == "a" {
        return "A"
    }
    if x == "b" {
        return "B"
    }

    return "unknown"
}
```