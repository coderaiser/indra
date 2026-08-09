# apply-dedent

> Remove unnecessary `dedent.Dedent(...)` wrappers

## Example

### ❌ Incorrect

```go
func f() []byte {
    return []byte(dedent.Dedent(`
        [ignore]
        patterns = ["vendor/**"]
        `))
}
```

### ✅ Correct

```go
func f() []byte {
    return []byte(`
        [ignore]
        patterns = ["vendor/**"]
        `)
}
```
