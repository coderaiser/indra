# conditions

Rules that apply to any assertion library. Nothing here is tape-specific, so
the group is enabled for all files by default — no glob restriction needed.

## Rules

- ✅ [remove-useless-condition](#remove-useless-condition)

## remove-useless-condition

Removes a useless condition inside `Ok()` / `NotOk()`.

### ❌ Incorrect

```go
package fixture

func f() {
	t.Ok(err != nil)
	t.NotOk(err == nil)
}
```

### ✅ Correct

```go
package fixture

func f() {
	t.Ok(err)
	t.NotOk(err)
}
```