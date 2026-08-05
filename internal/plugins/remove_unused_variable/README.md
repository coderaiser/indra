# remove-unused-variable

> Remove unused variables

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

## Disable

    [rules]
    "remove-unused-variable" = "off"
