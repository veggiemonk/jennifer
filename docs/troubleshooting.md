# Troubleshooting Jennifer

## Error handling basics

Jennifer's `File.Render()` and `File.Save()` return errors instead of panicking. Always check them:

```go
f := jen.NewFile("main")
// ... build code ...

if err := f.Save("output.go"); err != nil {
    log.Fatal(err)
}
```

The `%#v` / `GoString()` shorthand is convenient for tests but panics on error. Use `Render()` in production code.

## Reading error messages

Errors include the full path through the code tree. For example:

```
rendering statement item 3 (*jen.Group):
  rendering block item 0 (*jen.Statement):
    rendering statement item 2 (jen.token):
      unsupported type for Lit: struct { Name string } (value: {oops})
```

Read bottom-up: the root cause is the last line. The lines above trace which construct contains the problem — item 2 of a statement, inside item 0 of a block, inside item 3 of the outer statement.

## Common errors

### "unsupported type for Lit"

`Lit()` only supports Go's built-in types: `bool`, `string`, `int`, `float64`, `float32`, `int8`–`int64`, `uint`–`uint64`, `uintptr`, `complex64`, `complex128`.

```go
// Wrong — struct is not a built-in type
jen.Lit(MyStruct{Name: "foo"})

// Right — build the composite literal with Values
jen.ID("MyStruct").Values(jen.Dict{
    jen.ID("Name"): jen.Lit("foo"),
})
```

For `byte` and `rune` literals, use `LitByte()` and `LitRune()` instead of `Lit()`:

```go
jen.LitRune('x')    // renders: 'x'
jen.LitByte(0xab)   // renders: byte(0xab)
```

### "formatting generated source"

This means `go/format` rejected the generated code — it's syntactically invalid Go. The error includes the full numbered source so you can find the problem line.

To get structured error details, unwrap the `scanner.ErrorList`:

```go
import (
    "errors"
    "go/scanner"
)

err := f.Render(buf)
var se scanner.ErrorList
if errors.As(err, &se) {
    for _, e := range se {
        fmt.Printf("line %d, col %d: %s\n", e.Pos.Line, e.Pos.Column, e.Msg)
    }
}
```

### "Values: when using Dict, it must be the only item"

`Dict` must be the sole argument to `Values()`:

```go
// Wrong — Dict mixed with other items
jen.Values(jen.Dict{jen.ID("a"): jen.Lit(1)}, jen.Lit(2))

// Right — Dict is the only item
jen.Values(jen.Dict{
    jen.ID("a"): jen.Lit(1),
})
```

## Debugging techniques

### Render individual fragments

You don't have to render the whole file to test a fragment. Render a single statement:

```go
s := jen.If(jen.ID("err").Op("!=").Nil()).Block(
    jen.Return(jen.Err()),
)
fmt.Printf("%#v\n", s)
// Output: if err != nil { return err }
```

### Use NoFormat to see raw output

If `go/format` rejects your output and the error is hard to read, disable formatting to see exactly what jennifer produced:

```go
f := jen.NewFile("main")
f.NoFormat = true
// ... build code ...
f.Render(os.Stdout) // raw, unformatted output
```

### Inspect imports

Import conflicts (two packages with the same name) are resolved automatically by appending a number. If you see unexpected aliases like `http1`, check for conflicts:

```go
f := jen.NewFile("main")
f.Func().ID("main").Params().Block(
    jen.Qual("net/http", "Get").Call(jen.Lit("/")),
    jen.Qual("github.com/someone/http", "Do").Call(),
)
fmt.Printf("%#v\n", f)
// One of the imports will be aliased: http1 "github.com/someone/http"
```

To control aliases explicitly:

```go
f.ImportName("github.com/someone/http", "xhttp")
// or
f.ImportAlias("github.com/someone/http", "xhttp")
```

`ImportName` omits the alias when the name matches the package. `ImportAlias` always renders the alias.
