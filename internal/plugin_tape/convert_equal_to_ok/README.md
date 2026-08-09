# convert-equal-to-ok

> Convert `Equal(x, true)` to `Ok(x)`

## Example

### ❌ Incorrect

```go
func f(t *Test.T) {
    t.Equal(ok, true)
}
```

### ✅ Correct

```go
func f(t *Test.T) {
    t.Ok(ok)
}
```
