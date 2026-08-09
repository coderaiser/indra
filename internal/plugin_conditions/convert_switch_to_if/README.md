# convert-switch-to-if

> Use `if` instead of `switch` for one-value-per-case switches that end each
> case with a `return`

## Example

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

## Configuration

```toml
[rules]
"conditions/convert-switch-to-if" = "off"
```
