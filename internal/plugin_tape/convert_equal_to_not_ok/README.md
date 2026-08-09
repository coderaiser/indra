# convert-equal-to-not-ok

> Convert `Equal(x, nil/false)` to `NotOk(x)`

## Example

### ❌ Incorrect

```go
func f(t *Test.T) {
    t.Equal(err, nil)
}
```

### ✅ Correct

```go
func f(t *Test.T) {
    t.NotOk(err)
}
```
