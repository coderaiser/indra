# remove-skip

> Replace `Test.Skip` with `Test`

## Example

### ❌ Incorrect

```go
Test.Skip(t, "foo: bar", func(t *Test.T) {
    t.Equal(1, 1)
    t.End()
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
    "tape/remove-skip" = "off"
