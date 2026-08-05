# extract-result-from-assertion

> Extract inline expressions from assertions

## Example

### ❌ Incorrect

```go
t.DeepEqual(parse("input"), []Block{})
```

### ✅ Correct

```go
result := parse("input")
expected := []Block{}
t.DeepEqual(result, expected)
```

## Disable

```toml
[rules]
"tape/extract-result-from-assertion" = "off"
```
