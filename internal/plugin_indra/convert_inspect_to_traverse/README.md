# convert-inspect-to-traverse

> Convert `ast.Inspect` to `path.Traverse`

## Example

### ❌ Incorrect

```go
func Traverse() Traverser {
    return Traverser{"*ast.ReturnStmt": func(p Path, push func(Path)) {
        ast.Inspect(p.Node, func(n ast.Node) bool {
            return true
        })
        push(p)
    }}
}
```

### ✅ Correct

```go
func Traverse() Traverser {
    return Traverser{"*ast.ReturnStmt": func(p Path, push func(Path)) {
        p.Traverse(map[string]func(Path){
            "*ast.ReturnStmt": func(cp Path) {
                push(cp)
            },
        })
    }}
}
```
