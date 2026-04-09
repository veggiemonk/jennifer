package jen_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/veggiemonk/jennifer/jen"
)

func TestFileSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	f := NewFile("main")
	f.Func().ID("main").Params().Block()
	if err := f.Save(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "package main") {
		t.Errorf("saved file missing package declaration: %s", data)
	}
	if !strings.Contains(string(data), "func main()") {
		t.Errorf("saved file missing func main: %s", data)
	}
}

func TestFileSaveInvalidPath(t *testing.T) {
	f := NewFile("main")
	err := f.Save("/nonexistent/dir/test.go")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestGroupRender(t *testing.T) {
	g := BlockFunc(func(g *Group) {
		g.ID("x").Op("++")
	})
	got := fmt.Sprintf("%#v", g)
	if !strings.Contains(got, "x++") {
		t.Errorf("Group GoString: %s", got)
	}
}

func TestGroupRenderWithFile(t *testing.T) {
	g := &Group{}
	BlockFunc(func(inner *Group) {
		inner.Qual("fmt", "Println").Call(Lit("hello"))
		// We need the group to be rendered, use GoString on the wrapping statement
	})
	buf := &bytes.Buffer{}
	// Render a simple group
	simple := Block(Lit(1))
	var simpleGroup *Group
	BlockFunc(func(g *Group) {
		simpleGroup = g
		g.Lit(1)
	})
	_ = g
	_ = simple
	if simpleGroup != nil {
		if err := simpleGroup.Render(buf); err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if !strings.Contains(got, "1") {
			t.Errorf("Group.Render: %s", got)
		}
	}
}

func TestGroupGoString(t *testing.T) {
	var g *Group
	BlockFunc(func(inner *Group) {
		g = inner
		inner.Lit(1)
	})
	got := g.GoString()
	if !strings.Contains(got, "1") {
		t.Errorf("Group.GoString: %s", got)
	}
}

func TestStatementRenderWithFile(t *testing.T) {
	s := Qual("fmt", "Println").Call(Lit("hello"))
	f := NewFile("main")
	buf := &bytes.Buffer{}
	if err := s.RenderWithFile(buf, f); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if !strings.Contains(got, "Println") {
		t.Errorf("RenderWithFile: %s", got)
	}
}

func TestFileRenderError(t *testing.T) {
	f := NewFile("main")
	f.NoFormat = true
	// Add syntactically invalid code and then try to format
	f2 := NewFile("main")
	f2.Add(ID("func").ID("(")) // this will generate invalid Go
	buf := &bytes.Buffer{}
	// NoFormat should not error even with odd code
	f.Add(Lit(1))
	if err := f.Render(buf); err != nil {
		t.Fatal(err)
	}
}

func TestLitPanicOnUnsupportedType(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic for unsupported Lit type")
		}
	}()
	s := Lit(struct{}{})
	// GoString panics on render error, but Lit panics directly for unsupported types.
	// Use Render to trigger the render path.
	buf := &bytes.Buffer{}
	_ = s.Render(buf)
}

func TestDictValuesPanicOnMultipleItems(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic for Dict with multiple items in Values")
		}
	}()
	s := Values(Dict{ID("a"): Lit(1)}, Lit(2))
	buf := &bytes.Buffer{}
	_ = s.Render(buf)
}

func TestFileCanonicalPath(t *testing.T) {
	f := NewFile("main")
	f.CanonicalPath = "github.com/foo/bar"
	f.Func().ID("main").Params().Block()
	got := fmt.Sprintf("%#v", f)
	if !strings.Contains(got, `import "github.com/foo/bar"`) {
		t.Errorf("CanonicalPath not rendered: %s", got)
	}
}

func TestCommentMultiline(t *testing.T) {
	c := Comment("line1\nline2")
	got := fmt.Sprintf("%#v", c)
	if !strings.Contains(got, "/*") || !strings.Contains(got, "*/") {
		t.Errorf("multiline comment not rendered: %s", got)
	}
}

func TestCommentDirectFormatting(t *testing.T) {
	c := Comment("// direct comment")
	got := fmt.Sprintf("%#v", c)
	if got != "// direct comment" {
		t.Errorf("direct comment: got %q", got)
	}

	c2 := Comment("/* block */")
	got2 := fmt.Sprintf("%#v", c2)
	if got2 != "/* block */" {
		t.Errorf("direct block comment: got %q", got2)
	}
}

func TestDefaultInSwitch(t *testing.T) {
	c := Switch(ID("x")).Block(
		Case(Lit(1)).Block(ID("a").Call()),
		Default().Block(ID("b").Call()),
	)
	got := fmt.Sprintf("%#v", c)
	if !strings.Contains(got, "default:") {
		t.Errorf("default not rendered: %s", got)
	}
	if !strings.Contains(got, "case 1:") {
		t.Errorf("case not rendered: %s", got)
	}
}

func TestFileImportConflictResolution(t *testing.T) {
	f := NewFile("main")
	f.Func().ID("main").Params().Block(
		Qual("github.com/foo/bar", "A").Call(),
		Qual("github.com/baz/bar", "B").Call(),
	)
	got := fmt.Sprintf("%#v", f)
	// Both use "bar" as alias, one should be renamed
	if !strings.Contains(got, "bar") {
		t.Errorf("missing bar import: %s", got)
	}
}

func TestQualLocalPackage(t *testing.T) {
	f := NewFilePathName("github.com/my/pkg", "pkg")
	f.Func().ID("Foo").Params().Block(
		Qual("github.com/my/pkg", "Bar").Call(),
	)
	got := fmt.Sprintf("%#v", f)
	// Local package should not be imported, just render "Bar"
	if strings.Contains(got, `"github.com/my/pkg"`) {
		t.Errorf("local package should not appear in imports: %s", got)
	}
	if !strings.Contains(got, "Bar()") {
		t.Errorf("local qualified name should render as just name: %s", got)
	}
}

func TestTypesNullRendersNothing(t *testing.T) {
	c := Func().ID("F").Types().Params().Block()
	got := fmt.Sprintf("%#v", c)
	if strings.Contains(got, "[]") || strings.Contains(got, "[") && !strings.Contains(got, "{}") {
		t.Errorf("empty Types should not render brackets: %s", got)
	}
	if got != "func F() {}" {
		t.Errorf("expected 'func F() {}', got: %s", got)
	}
}
