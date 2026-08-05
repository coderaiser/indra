# convert-equal-to-deep-equal

> Use `DeepEqual` for slice args

## Example

### ❌ Incorrect

```go
t.Equal(result, []Block{})
```

### ✅ Correct

```go
t.DeepEqual(result, []Block{})
```

## Disable

    [rules]
    "tape/convert-equal-to-deep-equal" = "off"
