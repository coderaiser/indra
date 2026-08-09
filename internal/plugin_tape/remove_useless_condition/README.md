# remove-useless-condition

> Remove useless condition inside `Ok`/`NotOk`

## Example

### ❌ Incorrect

```go
func f(t *Test.T) {
    t.Ok(x != nil)
    t.NotOk(x == nil)
    t.Ok(x != false)
    t.NotOk(x == false)
}
```

### ✅ Correct

```go
func f(t *Test.T) {
    t.Ok(x)
    t.NotOk(x)
    t.Ok(x)
    t.NotOk(x)
}
```
