# convert-ok-to-not-ok

> Convert `Ok(err == nil)` and `Ok(!err)` to `NotOk(err)`

## Example

### ❌ Incorrect

```go
func f(t *Test.T) {
    t.Ok(err == nil)
    t.Ok(!err)
}
```

### ✅ Correct

```go
func f(t *Test.T) {
    t.NotOk(err)
    t.NotOk(err)
}
```
