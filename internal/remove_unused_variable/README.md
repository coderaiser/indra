# remove-unused-variable

> Remove unused variables

## Rule

```toml
[rules]
"remove-unused-variable" = "on"
```

### ❌ Incorrect

```go
//go:build ignore

package fixture

func f() int {
	x := 1
	return 0
}
```

### ✅ Correct

```go
//go:build ignore

package fixture

func f() int {

	return 0
}
```
