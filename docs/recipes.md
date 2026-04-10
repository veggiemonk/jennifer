# Jennifer Recipes

Patterns and techniques for common code generation tasks.

## Build constraints (build tags)

Use `HeaderComment` — it renders above the package clause with a blank line
separator, which is exactly where `go vet` expects build constraints:

```go
f := jen.NewFile("main")
f.HeaderComment("//go:build !windows")
```

Output:

```go
//go:build !windows

package main
```

Multiple constraints work too:

```go
f.HeaderComment("//go:build (linux && amd64) || darwin")
```

## Custom literal formats (hex, binary, octal)

`Lit()` always renders integers in decimal. If you need hex, binary, or octal notation, use `Op()` with a formatted string:

```go
jen.Op("0x80")           // renders: 0x80
jen.Op("0b1111_0000")    // renders: 0b1111_0000
jen.Op("0o755")          // renders: 0o755
```

For convenience, write a small helper in your own code:

```go
func LitHex(v int) *jen.Statement {
    return jen.Op(fmt.Sprintf("0x%X", v))
}
```

This keeps jennifer focused on correctness while you control the formatting.

## Build incrementally

When generating complex code, build and test one construct at a time:

```go
// Step 1: get the struct right
structCode := jen.Type().ID("Config").Struct(
    jen.ID("Host").String(),
    jen.ID("Port").Int(),
)
fmt.Printf("%#v\n\n", structCode)

// Step 2: add a constructor
newFunc := jen.Func().ID("NewConfig").Params(/* ... */).Op("*").ID("Config").Block(/* ... */)
fmt.Printf("%#v\n\n", newFunc)

// Step 3: combine into a file
f := jen.NewFile("config")
f.Add(structCode)
f.Add(newFunc)
```

## Generate code from reflect.Type

Jennifer doesn't have built-in reflection support, but you can build a helper
that converts a `reflect.Type` into jennifer code. This is useful when generating
code based on runtime type information (e.g., serializers, adapters, test fixtures).

```go
func QualFromType(tp reflect.Type) *jen.Statement {
    // Named type with a package path
    if tp.PkgPath() != "" && tp.Name() != "" {
        return jen.Qual(tp.PkgPath(), tp.Name())
    }
    // Built-in type (int, string, error, etc.)
    if tp.Name() != "" {
        return jen.ID(tp.Name())
    }
    // Composite types
    switch tp.Kind() {
    case reflect.Pointer:
        return jen.Op("*").Add(QualFromType(tp.Elem()))
    case reflect.Slice:
        return jen.Index().Add(QualFromType(tp.Elem()))
    case reflect.Array:
        return jen.Index(jen.Lit(tp.Len())).Add(QualFromType(tp.Elem()))
    case reflect.Map:
        return jen.Map(QualFromType(tp.Key())).Add(QualFromType(tp.Elem()))
    case reflect.Chan:
        switch tp.ChanDir() {
        case reflect.RecvDir:
            return jen.Op("<-").Chan().Add(QualFromType(tp.Elem()))
        case reflect.SendDir:
            return jen.Chan().Op("<-").Add(QualFromType(tp.Elem()))
        default:
            return jen.Chan().Add(QualFromType(tp.Elem()))
        }
    case reflect.Func:
        return jen.Func().ParamsFunc(func(g *jen.Group) {
            for i := range tp.NumIn() {
                g.Add(QualFromType(tp.In(i)))
            }
        }).ParamsFunc(func(g *jen.Group) {
            for i := range tp.NumOut() {
                g.Add(QualFromType(tp.Out(i)))
            }
        })
    default:
        panic(fmt.Errorf("unsupported reflect.Kind: %v", tp.Kind()))
    }
}
```

Adapt this to your needs — add struct tag handling, skip unexported fields,
handle generics, etc. It's intentionally not part of jennifer because every
project's requirements are slightly different.
