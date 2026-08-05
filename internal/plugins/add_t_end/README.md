# add-t-end

> Add missing `t.End()`

## Example

### ❌ Incorrect

```go
Test(t, "foo: bar", func(t *Test.T) {
    t.Equal(1, 1)
})
```

### ✅ Correct

```go
Test(t, "foo: bar", func(t *Test.T) {
    t.Equal(1, 1)
    t.End()
})
```

## Disable

    [rules]
    "tape/add-t-end" = "off"
