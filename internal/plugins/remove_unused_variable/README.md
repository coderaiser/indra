# remove-unused-variable

Remove unused variables

## Rule

```toml
[rules]
"remove-unused-variable" = "on"
```

## Example

### ❌ Incorrect

```go
x := compute()
return 0
```

### ✅ Correct

```go
return compute()
```
