package jen_test

import (
	"fmt"
	"go/format"
	"strings"
	"testing"

	. "github.com/veggiemonk/jennifer/jen"
)

var o1 = Options{
	Close:     ")",
	Multi:     true,
	Open:      "(",
	Separator: ",",
}

var o2 = Options{
	Close:     "",
	Multi:     false,
	Open:      "",
	Separator: ",",
}

type tc struct {
	desc          string
	code          Code
	expect        string
	expectImports map[string]string
}

var cases = []tc{
	{
		desc: `union_group`,
		code: Type().ID("A").InterfaceFunc(func(g *Group) {
			g.Union(ID("A"), ID("B"))
		}),
		expect: `type A interface{
			A|B
		}`,
	},
	{
		desc: `union_group_func`,
		code: Type().ID("A").InterfaceFunc(func(g1 *Group) {
			g1.UnionFunc(func(g2 *Group) {
				g2.ID("A")
				g2.ID("B")
			})
		}),
		expect: `type A interface{
			A|B
		}`,
	},
	{
		desc: `union`,
		code: Type().ID("A").Interface(Union(ID("A"), ID("B"))),
		expect: `type A interface{
			A|B
		}`,
	},
	{
		desc: `unionFunc`,
		code: Type().ID("A").Interface(UnionFunc(func(g *Group) {
			g.ID("A")
			g.ID("B")
		})),
		expect: `type A interface{
			A|B
		}`,
	},
	{
		desc:   `types1`,
		code:   Func().ID("A").Types(ID("K").Comparable(), ID("V").Any()).Params(),
		expect: `func A[K comparable, V any]()`,
	},
	{
		desc:   `types2`,
		code:   Func().ID("A").Types(ID("T1"), ID("T2").Any()).Params(),
		expect: `func A[T1, T2 any]()`,
	},
	{
		desc:   `types func`,
		code:   Func().ID("A").Add(Types(ID("T1"), ID("T2").Any())).Params(),
		expect: `func A[T1, T2 any]()`,
	},
	{
		desc:   `scientific notation`,
		code:   Lit(1e3),
		expect: `1000.0`,
	},
	{
		desc:   `big float`,
		code:   Lit(1000000.0),
		expect: `1e+06`,
	},
	{
		desc:   `lit float whole numbers`,
		code:   Index().Float64().Values(Lit(-10.0), Lit(-2.0), Lit(-1.0), Lit(0.0), Lit(1.0), Lit(2.0), Lit(10.0)),
		expect: "[]float64{-10.0, -2.0, -1.0, 0.0, 1.0, 2.0, 10.0}",
	},
	{
		desc: `custom func group`,
		code: ListFunc(func(g *Group) {
			g.CustomFunc(o2, func(g *Group) {
				g.ID("a")
				g.ID("b")
				g.ID("c")
			})
		}).Op("=").ID("foo").Call(),
		expect: `a, b, c = foo()`,
	},
	{
		desc:   `custom group`,
		code:   ListFunc(func(g *Group) { g.Custom(o2, ID("a"), ID("b"), ID("c")) }).Op("=").ID("foo").Call(),
		expect: `a, b, c = foo()`,
	},
	{
		desc: `custom function`,
		code: ID("foo").Add(Custom(o1, Lit("a"), Lit("b"), Lit("c"))),
		expect: `foo(
			"a",
			"b",
			"c",
		)`,
	},
	{
		desc: `line statement`,
		code: Block(Lit(1).Line(), Lit(2)),
		expect: `{
		1

		2
		}`,
	},
	{
		desc: `line func`,
		code: Block(Lit(1), Line(), Lit(2)),
		expect: `{
		1

		2
		}`,
	},
	{
		desc: `line group`,
		code: BlockFunc(func(g *Group) {
			g.ID("a")
			g.Line()
			g.ID("b")
		}),
		expect: `{
		a

		b
		}`,
	},
	{
		desc: `op group`,
		code: BlockFunc(func(g *Group) {
			g.Op("*").ID("a")
		}),
		expect: `{*a}`,
	},
	{
		desc: `empty group`,
		code: BlockFunc(func(g *Group) {
			g.Empty()
		}),
		expect: `{

		}`,
	},
	{
		desc: `null group`,
		code: BlockFunc(func(g *Group) {
			g.Null()
		}),
		expect: `{}`,
	},
	{
		desc:   `tag no backquote`,
		code:   Tag(map[string]string{"a": "`b`"}),
		expect: "\"a:\\\"`b`\\\"\"",
	},
	{
		desc:   `tag null`,
		code:   Tag(map[string]string{}),
		expect: ``,
	},
	{
		desc: `litbytefunc group`,
		code: BlockFunc(func(g *Group) {
			g.LitByteFunc(func() byte { return byte(0xab) })
		}),
		expect: `{byte(0xab)}`,
	},
	{
		desc: `litbyte group`,
		code: BlockFunc(func(g *Group) {
			g.LitByte(byte(0xab))
		}),
		expect: `{byte(0xab)}`,
	},
	{
		desc: `litrunefunc group`,
		code: BlockFunc(func(g *Group) {
			g.LitRuneFunc(func() rune { return 'a' })
		}),
		expect: `{'a'}`,
	},
	{
		desc: `litrune group`,
		code: BlockFunc(func(g *Group) {
			g.LitRune('a')
		}),
		expect: `{'a'}`,
	},
	{
		desc: `litfunc group`,
		code: BlockFunc(func(g *Group) {
			g.LitFunc(func() any {
				return 1 + 1
			})
		}),
		expect: `{2}`,
	},
	{
		desc: `litfunc func`,
		code: LitFunc(func() any {
			return 1 + 1
		}),
		expect: `2`,
	},
	{
		desc:   `group all null`,
		code:   List(Null(), Null()),
		expect: ``,
	},
	{
		desc:   `do group`,
		code:   BlockFunc(func(g *Group) { g.Do(func(s *Statement) { s.Lit(1) }) }),
		expect: `{1}`,
	},
	{
		desc:   `do func`,
		code:   Do(func(s *Statement) { s.Lit(1) }),
		expect: `1`,
	},
	{
		desc:   `dict empty`,
		code:   Values(Dict{}),
		expect: `{}`,
	},
	{
		desc:   `dict null`,
		code:   Values(Dict{Null(): Null()}),
		expect: `{}`,
	},
	{
		desc: `commentf group`,
		code: BlockFunc(func(g *Group) { g.Commentf("%d", 1) }),
		expect: `{
		// 1
		}`,
	},
	{
		desc:   `commentf func`,
		code:   Commentf("%d", 1),
		expect: `// 1`,
	},
	{
		desc:   `add func`,
		code:   Add(Lit(1)),
		expect: `1`,
	},
	{
		desc:   `add group`,
		code:   BlockFunc(func(g *Group) { g.Add(Lit(1)) }),
		expect: `{1}`,
	},
	{
		desc:   `empty block`,
		code:   Block(),
		expect: `{}`,
	},
	{
		desc:   `string literal`,
		code:   Lit("a"),
		expect: `"a"`,
	},
	{
		desc:   `int literal`,
		code:   Lit(1),
		expect: `1`,
	},
	{
		desc:   `simple id`,
		code:   ID("a"),
		expect: `a`,
	},
	{
		desc:   `foreign id`,
		code:   Qual("x.y/z", "a"),
		expect: `z.a`,
		expectImports: map[string]string{
			"x.y/z": "z",
		},
	},
	{
		desc:   `var decl`,
		code:   Var().ID("a").Op("=").Lit("b"),
		expect: `var a = "b"`,
	},
	{
		desc:   `short var decl`,
		code:   ID("a").Op(":=").Lit("b"),
		expect: `a := "b"`,
	},
	{
		desc:   `simple if`,
		code:   If(ID("a").Op("==").Lit("b")).Block(),
		expect: `if a == "b" {}`,
	},
	{
		desc: `simple if with body`,
		code: If(ID("a").Op("==").Lit("b")).Block(
			ID("a").Op("++"),
		),
		expect: `if a == "b" { a++ }`,
	},
	{
		desc:   `pointer`,
		code:   Op("*").ID("a"),
		expect: `*a`,
	},
	{
		desc:   `address`,
		code:   Op("&").ID("a"),
		expect: `&a`,
	},
	{
		desc:   `simple call`,
		code:   ID("a").Call(Lit("b"), Lit("c")),
		expect: `a("b", "c")`,
	},
	{
		desc:   `call fmt.Sprintf`,
		code:   Qual("fmt", "Sprintf").Call(Lit("b"), ID("c")),
		expect: `fmt.Sprintf("b", c)`,
	},
	{
		desc:   `slices`,
		code:   ID("a").Index(Lit(1), Empty()),
		expect: `a[1:]`,
	},
	{
		desc:   `return`,
		code:   Return(ID("a")),
		expect: `return a`,
	},
	{
		desc:   `double return`,
		code:   Return(ID("a"), ID("b")),
		expect: `return a, b`,
	},
	{
		desc: `func`,
		code: Func().ID("a").Params(
			ID("a").String(),
		).Block(
			Return(ID("a")),
		),
		expect: `func a(a string){
			return a
		}`,
	},
	{
		desc:   `built in func`,
		code:   New(ID("a")),
		expect: `new(a)`,
	},
	{
		desc:   `multip`,
		code:   ID("a").Op("*").ID("b"),
		expect: `a * b`,
	},
	{
		desc:   `multip ptr`,
		code:   ID("a").Op("*").Op("*").ID("b"),
		expect: `a * *b`,
	},
	{
		desc:   `field`,
		code:   ID("a").Dot("b"),
		expect: `a.b`,
	},
	{
		desc:   `method`,
		code:   ID("a").Dot("b").Call(ID("c"), ID("d")),
		expect: `a.b(c, d)`,
	},
	{
		desc: `if else`,
		code: If(ID("a").Op("==").Lit(1)).Block(
			ID("b").Op("=").Lit(1),
		).Else().If(ID("a").Op("==").Lit(2)).Block(
			ID("b").Op("=").Lit(2),
		).Else().Block(
			ID("b").Op("=").Lit(3),
		),
		expect: `if a == 1 { b = 1 } else if a == 2 { b = 2 } else { b = 3 }`,
	},
	{
		desc:   `literal array`,
		code:   Index().String().Values(Lit("a"), Lit("b")),
		expect: `[]string{"a", "b"}`,
	},
	{
		desc:   `comment`,
		code:   Comment("a"),
		expect: `// a`,
	},
	{
		desc:   `null`,
		code:   ID("a").Params(ID("b"), Null(), ID("c")),
		expect: `a(b, c)`,
	},
	{
		desc: `map literal single`,
		code: ID("a").Values(Dict{
			ID("b"): ID("c"),
		}),
		expect: `a{b: c}`,
	},
	{
		desc: `map literal null`,
		code: ID("a").Values(Dict{
			Null():  ID("c"),
			ID("b"): Null(),
			ID("b"): ID("c"),
		}),
		expect: `a{b: c}`,
	},
	{
		desc: `map literal multiple`,
		code: ID("a").Values(Dict{
			ID("b"): ID("c"),
			ID("d"): ID("e"),
		}),
		expect: `a{
			b: c,
			d: e,
		}`,
	},
	{
		desc: `map literal func single`,
		code: ID("a").Values(DictFunc(func(d Dict) {
			d[ID("b")] = ID("c")
		})),
		expect: `a{b: c}`,
	},
	{
		desc: `map literal func single null`,
		code: ID("a").Values(DictFunc(func(d Dict) {
			d[Null()] = ID("c")
			d[ID("b")] = Null()
			d[ID("b")] = ID("c")
		})),
		expect: `a{b: c}`,
	},
	{
		desc: `map literal func multiple`,
		code: ID("a").Values(DictFunc(func(d Dict) {
			d[ID("b")] = ID("c")
			d[ID("d")] = ID("e")
		})),
		expect: `a{
			b: c,
			d: e,
		}`,
	},
	{
		desc: `literal func`,
		code: ID("a").Op(":=").LitFunc(func() any {
			return "b"
		}),
		expect: `a := "b"`,
	},
	{
		desc:   `dot`,
		code:   ID("a").Dot("b").Dot("c"),
		expect: `a.b.c`,
	},
	{
		desc:   `do`,
		code:   ID("a").Do(func(s *Statement) { s.Dot("b") }),
		expect: `a.b`,
	},
	{
		desc:   `tags should be ordered`,
		code:   Tag(map[string]string{"z": "1", "a": "2"}),
		expect: "`a:\"2\" z:\"1\"`",
	},
	{
		desc: `dict should be ordered`,
		code: Map(String()).Int().Values(Dict{ID("z"): Lit(1), ID("a"): Lit(2)}),
		expect: `map[string]int{
		a:2,
		z:1,
		}`,
	},
	// switch/case
	{
		desc: `switch case block`,
		code: Switch(ID("a")).Block(
			Case(Lit(1)).Block(
				Var().ID("i").Int(),
				Var().ID("j").Int(),
			),
		),
		expect: `switch a {
		case 1:
			var i int
			var j int
		}`,
	},
	{
		desc: `switch default`,
		code: Switch(ID("a")).Block(
			Case(Lit(1)).Block(ID("a").Op("++")),
			Default().Block(ID("b").Op("++")),
		),
		expect: `switch a {
		case 1:
			a++
		default:
			b++
		}`,
	},
}

func TestJen(t *testing.T) {
	caseTester(t, cases)
}

func caseTester(t *testing.T, cases []tc) {
	for i, c := range cases {
		rendered := fmt.Sprintf("%#v", c.code)

		expected, err := format.Source([]byte(c.expect))
		if err != nil {
			panic(fmt.Sprintf("Error formatting expected source in test case %d. Description: %s\nError:\n%s", i, c.desc, err))
		}

		if strings.TrimSpace(rendered) != strings.TrimSpace(string(expected)) {
			t.Errorf("Test case %d failed. Description: %s\nExpected:\n%s\nOutput:\n%s", i, c.desc, expected, rendered)
		}

		if c.expectImports != nil {
			f := NewFile("test")
			f.Add(c.code)
			output := fmt.Sprintf("%#v", f)
			for path, alias := range c.expectImports {
				if !strings.Contains(output, alias+".") && !strings.Contains(output, `"`+path+`"`) {
					t.Errorf("Test case %d failed. Description: %s\nExpected import %s (%s) not found in:\n%s", i, c.desc, path, alias, output)
				}
			}
		}
	}
}

func TestNilStatement(t *testing.T) {
	var s *Statement
	c := Func().ID("a").Params(s)
	got := fmt.Sprintf("%#v", c)
	expect := "func a()"
	if got != expect {
		t.Fatalf("Got: %s, expect: %s", got, expect)
	}
}

func TestNilGroup(t *testing.T) {
	var g *Group
	c := Func().ID("a").Params(g)
	got := fmt.Sprintf("%#v", c)
	expect := "func a()"
	if got != expect {
		t.Fatalf("Got: %s, expect: %s", got, expect)
	}
}

func TestBlockAfterCaseRenderIdempotent(t *testing.T) {
	c := Switch(ID("x")).Block(
		Case(Lit(1)).Block(ID("a").Call()),
	)
	first := fmt.Sprintf("%#v", c)
	second := fmt.Sprintf("%#v", c)
	if first != second {
		t.Errorf("Block-after-Case renders differ:\nfirst:  %s\nsecond: %s", first, second)
	}
}

func TestImportAliasForAnyAndComparable(t *testing.T) {
	f := NewFile("main")
	c := Qual("example.com/any", "Foo")
	got := fmt.Sprintf("%#v", f.Func().ID("f").Params().Block(Add(c)))
	if strings.Contains(got, "import any ") {
		t.Errorf("import alias should not be 'any', got:\n%s", got)
	}
}
