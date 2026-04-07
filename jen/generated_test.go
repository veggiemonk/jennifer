package jen_test

import (
	"testing"

	. "github.com/veggiemonk/jennifer/jen"
)

var gencases = []tc{
	{
		desc: `typesfunc group`,
		// Don't do this! ListFunc used to kludge Group.TypesFunc usage without
		// syntax error.
		code:   ID("a").ListFunc(func(lg *Group) { lg.TypesFunc(func(cg *Group) { cg.Lit(1) }) }),
		expect: `a[1]`,
	},
	{
		desc: `types group`,
		// Don't do this! ListFunc used to kludge Group.Types usage without
		// syntax error.
		code:   ID("a").ListFunc(func(lg *Group) { lg.Types(Lit(1)) }),
		expect: `a[1]`,
	},
	{
		desc: `bool group`,
		code: BlockFunc(func(g *Group) {
			g.Bool()
		}),
		expect: `{
		bool
		}`,
	},
	{
		desc:   `recover func`,
		code:   Recover(),
		expect: `recover()`,
	},
	{
		desc:   `recover statement`,
		code:   Null().Recover(),
		expect: `recover()`,
	},
	{
		desc: `recover group`,
		code: BlockFunc(func(g *Group) {
			g.Recover()
		}),
		expect: `{
		recover()
		}`,
	},
	{
		desc:   `real func`,
		code:   Real(ID("a")),
		expect: `real(a)`,
	},
	{
		desc:   `real statement`,
		code:   Null().Real(ID("a")),
		expect: `real(a)`,
	},
	{
		desc: `real group`,
		code: BlockFunc(func(g *Group) {
			g.Real(ID("a"))
		}),
		expect: `{
		real(a)
		}`,
	},
	{
		desc:   `printlnfunc func`,
		code:   PrintlnFunc(func(g *Group) { g.ID("a") }),
		expect: `println(a)`,
	},
	{
		desc:   `printlnfunc statement`,
		code:   Null().PrintlnFunc(func(g *Group) { g.ID("a") }),
		expect: `println(a)`,
	},
	{
		desc: `printlnfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.PrintlnFunc(func(pg *Group) { pg.ID("a") })
		}),
		expect: `{
		println(a)
		}`,
	},
	{
		desc:   `println func`,
		code:   Println(ID("a")),
		expect: `println(a)`,
	},
	{
		desc:   `println statement`,
		code:   Null().Println(ID("a")),
		expect: `println(a)`,
	},
	{
		desc: `println group`,
		code: BlockFunc(func(g *Group) {
			g.Println(ID("a"))
		}),
		expect: `{
		println(a)
		}`,
	},
	{
		desc:   `printfunc func`,
		code:   PrintFunc(func(g *Group) { g.ID("a") }),
		expect: `print(a)`,
	},
	{
		desc:   `printfunc statement`,
		code:   Null().PrintFunc(func(g *Group) { g.ID("a") }),
		expect: `print(a)`,
	},
	{
		desc: `printfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.PrintFunc(func(pg *Group) { pg.ID("a") })
		}),
		expect: `{
		print(a)
		}`,
	},
	{
		desc:   `print func`,
		code:   Print(ID("a")),
		expect: `print(a)`,
	},
	{
		desc:   `print statement`,
		code:   Null().Print(ID("a")),
		expect: `print(a)`,
	},
	{
		desc: `print group`,
		code: BlockFunc(func(g *Group) {
			g.Print(ID("a"))
		}),
		expect: `{
		print(a)
		}`,
	},
	{
		desc:   `panic func`,
		code:   Panic(ID("a")),
		expect: `panic(a)`,
	},
	{
		desc:   `panic statement`,
		code:   Null().Panic(ID("a")),
		expect: `panic(a)`,
	},
	{
		desc: `panic group`,
		code: BlockFunc(func(g *Group) {
			g.Panic(ID("a"))
		}),
		expect: `{
		panic(a)
		}`,
	},
	{
		desc:   `new func`,
		code:   New(ID("a")),
		expect: `new(a)`,
	},
	{
		desc:   `new statement`,
		code:   ID("a").Op(":=").New(ID("a")),
		expect: `a := new(a)`,
	},
	{
		desc: `new group`,
		code: BlockFunc(func(g *Group) {
			g.New(ID("a"))
		}),
		expect: `{
		new(a)
		}`,
	},
	{
		desc:   `make func`,
		code:   Make(ID("a")),
		expect: `make(a)`,
	},
	{
		desc:   `make statement`,
		code:   ID("a").Op(":=").Make(ID("a")),
		expect: `a := make(a)`,
	},
	{
		desc: `make group`,
		code: BlockFunc(func(g *Group) {
			g.Make(ID("a"))
		}),
		expect: `{
		make(a)
		}`,
	},
	{
		desc:   `len func`,
		code:   Len(ID("a")),
		expect: `len(a)`,
	},
	{
		desc:   `len statement`,
		code:   ID("a").Op(":=").Len(ID("a")),
		expect: `a := len(a)`,
	},
	{
		desc: `len group`,
		code: BlockFunc(func(g *Group) {
			g.Len(ID("a"))
		}),
		expect: `{
		len(a)
		}`,
	},
	{
		desc:   `imag func`,
		code:   Imag(ID("a")),
		expect: `imag(a)`,
	},
	{
		desc:   `imag statement`,
		code:   ID("a").Op(":=").Imag(ID("a")),
		expect: `a := imag(a)`,
	},
	{
		desc: `imag group`,
		code: BlockFunc(func(g *Group) {
			g.Imag(ID("a"))
		}),
		expect: `{
		imag(a)
		}`,
	},
	{
		desc:   `delete func`,
		code:   Delete(ID("a"), ID("b")),
		expect: `delete(a, b)`,
	},
	{
		desc:   `delete statement`,
		code:   Null().Delete(ID("a"), ID("b")),
		expect: `delete(a, b)`,
	},
	{
		desc: `delete group`,
		code: BlockFunc(func(g *Group) {
			g.Delete(ID("a"), ID("b"))
		}),
		expect: `{
		delete(a, b)
		}`,
	},
	{
		desc:   `copy func`,
		code:   Copy(ID("a"), ID("b")),
		expect: `copy(a, b)`,
	},
	{
		desc:   `copy statement`,
		code:   ID("a").Op(":=").Copy(ID("a"), ID("b")),
		expect: `a := copy(a, b)`,
	},
	{
		desc: `copy group`,
		code: BlockFunc(func(g *Group) {
			g.Copy(ID("a"), ID("b"))
		}),
		expect: `{
		copy(a, b)
		}`,
	},
	{
		desc:   `complex func`,
		code:   Complex(ID("a"), ID("b")),
		expect: `complex(a, b)`,
	},
	{
		desc:   `complex statement`,
		code:   ID("a").Op(":=").Complex(ID("a"), ID("b")),
		expect: `a := complex(a, b)`,
	},
	{
		desc: `complex group`,
		code: BlockFunc(func(g *Group) {
			g.Complex(ID("a"), ID("b"))
		}),
		expect: `{
		complex(a, b)
		}`,
	},
	{
		desc: `close group`,
		code: BlockFunc(func(g *Group) { g.Close(ID("a")) }),
		expect: `{
		close(a)
		}`,
	},
	{
		desc:   `cap func`,
		code:   Cap(ID("a")),
		expect: `cap(a)`,
	},
	{
		desc:   `cap statement`,
		code:   ID("a").Op(":=").Cap(ID("b")),
		expect: `a := cap(b)`,
	},
	{
		desc: `cap group`,
		code: BlockFunc(func(g *Group) {
			g.Cap(ID("a"))
		}),
		expect: `{
		cap(a)
		}`,
	},
	{
		desc: `append group`,
		code: BlockFunc(func(g *Group) {
			g.Append(ID("a"))
		}),
		expect: `{
		append(a)
		}`,
	},
	{
		desc:   `appendfunc statement`,
		code:   ID("a").Op("=").AppendFunc(func(ag *Group) { ag.ID("a") }),
		expect: `a = append(a)`,
	},
	{
		desc:   `appendfunc func`,
		code:   AppendFunc(func(ag *Group) { ag.ID("a") }),
		expect: `append(a)`,
	},
	{
		desc: `appendfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.AppendFunc(func(ag *Group) { ag.ID("a") })
		}),
		expect: `{
		append(a)
		}`,
	},
	{
		desc: `casefunc group`,
		code: Switch().BlockFunc(func(g *Group) {
			g.CaseFunc(func(g *Group) { g.ID("a") }).Block()
		}),
		expect: `switch {
		case a:
		}`,
	},
	{
		desc: `case group`,
		code: Switch().BlockFunc(func(g *Group) {
			g.Case(ID("a")).Block()
		}),
		expect: `switch {
		case a:
		}`,
	},
	{
		desc:   `structfunc statement`,
		code:   ID("a").Op(":=").StructFunc(func(g *Group) {}).Values(),
		expect: `a := struct{}{}`,
	},
	{
		desc: `structfunc group`,
		// Don't do this! ListFunc used to kludge Group.Struct usage
		// without syntax error.
		code:   ID("a").Op(":=").ListFunc(func(g *Group) { g.StructFunc(func(g *Group) {}) }).Values(),
		expect: `a := struct{}{}`,
	},
	{
		desc:   `structfunc func`,
		code:   ID("a").Op(":=").Add(StructFunc(func(g *Group) {})).Values(),
		expect: `a := struct{}{}`,
	},
	{
		desc: `struct group`,
		// Don't do this! ListFunc used to kludge Group.Struct usage
		// without syntax error.
		code:   ID("a").Op(":=").ListFunc(func(g *Group) { g.Struct() }).Values(),
		expect: `a := struct{}{}`,
	},
	{
		desc:   `struct func`,
		code:   ID("a").Op(":=").Add(Struct()).Values(),
		expect: `a := struct{}{}`,
	},
	{
		desc: `interfacefunc func`,
		code: ID("a").Assert(InterfaceFunc(func(g *Group) {
			g.ID("a").Call().Int()
			g.ID("b").Call().Int()
		})),
		expect: `a.(interface{
		a() int
		b() int
		})`,
	},
	{
		desc: `interfacefunc statement`,
		code: ID("a").Assert(Null().InterfaceFunc(func(g *Group) {
			g.ID("a").Call().Int()
			g.ID("b").Call().Int()
		})),
		expect: `a.(interface{
		a() int
		b() int
		})`,
	},
	{
		desc: `interfacefunc group`,
		// Don't do this! ListFunc used to kludge Group.InterfaceFunc usage
		// without syntax error.
		code: ID("a").Assert(ListFunc(func(lg *Group) {
			lg.InterfaceFunc(func(ig *Group) {
				ig.ID("a").Call().Int()
				ig.ID("b").Call().Int()
			})
		})),
		expect: `a.(interface{
		a() int
		b() int
		})`,
	},
	{
		desc:   `interface func`,
		code:   Interface().Parens(ID("a")),
		expect: `interface{}(a)`,
	},
	{
		desc: `interface group`,
		code: BlockFunc(func(g *Group) {
			g.Interface().Parens(ID("a"))
		}),
		expect: `{
		interface{}(a)
		}`,
	},
	{
		desc:   `interface statement`,
		code:   Null().Interface().Parens(ID("a")),
		expect: `interface{}(a)`,
	},
	{
		desc: `switchfunc func`,
		code: SwitchFunc(func(rg *Group) {
			rg.ID("a")
		}).Block(),
		expect: `switch a {}`,
	},
	{
		desc: `switchfunc statement`,
		code: Null().SwitchFunc(func(rg *Group) {
			rg.ID("a")
		}).Block(),
		expect: `switch a {
		}`,
	},
	{
		desc: `switchfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.SwitchFunc(func(rg *Group) {
				rg.ID("a")
			}).Block()
		}),
		expect: `{
			switch a {
			}
		}`,
	},
	{
		desc: `switch group`,
		code: BlockFunc(func(bg *Group) {
			bg.Switch().Block()
		}),
		expect: `{
			switch {
			}
		}`,
	},
	{
		desc: `forfunc func`,
		code: ForFunc(func(rg *Group) {
			rg.ID("a")
		}).Block(),
		expect: `for a {}`,
	},
	{
		desc: `forfunc statement`,
		code: Null().ForFunc(func(rg *Group) {
			rg.ID("a")
		}).Block(),
		expect: `for a {
		}`,
	},
	{
		desc: `forfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.ForFunc(func(rg *Group) {
				rg.ID("a")
			}).Block()
		}),
		expect: `{
			for a {
			}
		}`,
	},
	{
		desc: `for group`,
		code: BlockFunc(func(g *Group) {
			g.For(ID("a")).Block()
		}),
		expect: `{
		for a {}
		}`,
	},
	{
		desc: `returnfunc func`,
		code: ReturnFunc(func(rg *Group) {
			rg.Lit(1)
			rg.Lit(2)
		}),
		expect: `return 1, 2`,
	},
	{
		desc: `returnfunc statement`,
		code: Empty().ReturnFunc(func(rg *Group) {
			rg.Lit(1)
			rg.Lit(2)
		}),
		expect: `return 1, 2`,
	},
	{
		desc: `returnfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.ReturnFunc(func(rg *Group) {
				rg.Lit(1)
				rg.Lit(2)
			})
		}),
		expect: `{
		return 1, 2
		}`,
	},
	{
		desc: `return group`,
		code: BlockFunc(func(g *Group) {
			g.Return()
		}),
		expect: `{
		return
		}`,
	},
	{
		desc: `iffunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.IfFunc(func(ig *Group) {
				ig.ID("a")
			}).Block()
		}),
		expect: `{
		if a {} 
		}`,
	},
	{
		desc: `iffunc func`,
		code: IfFunc(func(ig *Group) {
			ig.ID("a")
		}).Block(),
		expect: `if a {}`,
	},
	{
		desc: `iffunc statement`,
		code: Null().IfFunc(func(ig *Group) {
			ig.ID("a")
		}).Block(),
		expect: `if a {}`,
	},
	{
		desc: `if group`,
		code: BlockFunc(func(g *Group) { g.If(ID("a")).Block() }),
		expect: `{
		if a {}
		}`,
	},
	{
		desc: `map group`,
		code: BlockFunc(func(g *Group) { g.Map(Int()).Int().Values(Dict{Lit(1): Lit(1)}) }),
		expect: `{
		map[int]int{1:1}
		}`,
	},
	{
		desc: `assert group`,
		// Don't do this! ListFunc used to kludge Group.Assert usage without
		// syntax error.
		code:   ID("a").ListFunc(func(g *Group) { g.Assert(ID("b")) }),
		expect: `a.(b)`,
	},
	{
		desc:   `assert func`,
		code:   ID("a").Add(Assert(ID("b"))),
		expect: `a.(b)`,
	},
	{
		desc: `paramsfunc group`,
		// Don't do this! ListFunc used to kludge Group.ParamsFunc usage without
		// syntax error.
		code:   ID("a").ListFunc(func(lg *Group) { lg.ParamsFunc(func(cg *Group) { cg.Lit(1) }) }),
		expect: `a(1)`,
	},
	{
		desc:   `paramsfunc func`,
		code:   ID("a").Add(ParamsFunc(func(g *Group) { g.Lit(1) })),
		expect: `a(1)`,
	},
	{
		desc:   `paramsfunc statement`,
		code:   ID("a").ParamsFunc(func(g *Group) { g.Lit(1) }),
		expect: `a(1)`,
	},
	{
		desc: `params group`,
		// Don't do this! ListFunc used to kludge Group.Params usage without
		// syntax error.
		code:   ID("a").ListFunc(func(g *Group) { g.Params(Lit(1)) }),
		expect: `a(1)`,
	},
	{
		desc:   `params func`,
		code:   ID("a").Add(Params(Lit(1))),
		expect: `a(1)`,
	},
	{
		desc: `callfunc group`,
		// Don't do this! ListFunc used to kludge Group.CallFunc usage without
		// syntax error.
		code:   ID("a").ListFunc(func(lg *Group) { lg.CallFunc(func(cg *Group) { cg.Lit(1) }) }),
		expect: `a(1)`,
	},
	{
		desc:   `callfunc func`,
		code:   ID("a").Add(CallFunc(func(g *Group) { g.Lit(1) })),
		expect: `a(1)`,
	},
	{
		desc: `call group`,
		// Don't do this! ListFunc used to kludge Group.Call usage without
		// syntax error.
		code:   ID("a").ListFunc(func(g *Group) { g.Call(Lit(1)) }),
		expect: `a(1)`,
	},
	{
		desc:   `call func`,
		code:   ID("a").Add(Call(Lit(1))),
		expect: `a(1)`,
	},
	{
		desc: `defsfunc statement`,
		code: Const().DefsFunc(func(g *Group) { g.ID("a").Op("=").Lit(1) }),
		expect: `const (
		a = 1
		)`,
	},
	{
		desc: `defsfunc func`,
		code: Const().Add(DefsFunc(func(g *Group) { g.ID("a").Op("=").Lit(1) })),
		expect: `const (
		a = 1
		)`,
	},
	{
		desc: `defsfunc group`,
		// Don't do this! ListFunc used to kludge Group.DefsFunc usage without
		// syntax error.
		code: Const().ListFunc(func(lg *Group) { lg.DefsFunc(func(dg *Group) { dg.ID("a").Op("=").Lit(1) }) }),
		expect: `const (
		a = 1
		)`,
	},
	{
		desc: `defs group`,
		// Don't do this! ListFunc used to kludge Group.Defs usage without
		// syntax error.
		code: Const().ListFunc(func(g *Group) { g.Defs(ID("a").Op("=").Lit(1)) }),
		expect: `const (
		a = 1
		)`,
	},
	{
		desc: `defs func`,
		code: Const().Add(Defs(ID("a").Op("=").Lit(1))),
		expect: `const (
		a = 1
		)`,
	},
	{
		desc:   `blockfunc group`,
		code:   BlockFunc(func(g *Group) { g.BlockFunc(func(g *Group) {}) }),
		expect: `{{}}`,
	},
	{
		desc:   `block group`,
		code:   BlockFunc(func(g *Group) { g.Block() }),
		expect: `{{}}`,
	},
	{
		desc:   `indexfunc group`,
		code:   BlockFunc(func(g *Group) { g.IndexFunc(func(g *Group) { g.Lit(1) }).Int().Values(Lit(1)) }),
		expect: `{[1]int{1}}`,
	},
	{
		desc:   `indexfunc statement`,
		code:   ID("a").IndexFunc(func(g *Group) { g.Lit(1) }),
		expect: `a[1]`,
	},
	{
		desc:   `indexfunc func`,
		code:   ID("a").Add(IndexFunc(func(g *Group) { g.Lit(1) })),
		expect: `a[1]`,
	},
	{
		desc:   `index group`,
		code:   BlockFunc(func(g *Group) { g.Index(Lit(1)).Int().Values(Lit(1)) }),
		expect: `{[1]int{1}}`,
	},
	{
		desc:   `index func`,
		code:   ID("a").Add(Index(Lit(1))),
		expect: `a[1]`,
	},
	{
		desc: `valuesfunc func`,
		code: ValuesFunc(func(vg *Group) {
			vg.Lit(1)
		}),
		expect: `{1}`,
	},
	{
		desc: `valuesfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.ValuesFunc(func(vg *Group) {
				vg.Lit(1)
			})
		}),
		expect: `{
		{1}
		}`,
	},
	{
		desc: `values group`,
		code: BlockFunc(func(g *Group) {
			g.Values(Lit(1))
		}),
		expect: `{
		{1}
		}`,
	},
	{
		desc: `listfunc statement`,
		code: Add(Null()).ListFunc(func(lg *Group) {
			lg.ID("a")
			lg.ID("b")
		}).Op("=").ID("c"),
		expect: `a, b = c`,
	},
	{
		desc: `listfunc func`,
		code: ListFunc(func(lg *Group) {
			lg.ID("a")
			lg.ID("b")
		}).Op("=").ID("c"),
		expect: `a, b = c`,
	},
	{
		desc: `listfunc group`,
		code: BlockFunc(func(bg *Group) {
			bg.ListFunc(func(lg *Group) {
				lg.ID("a")
				lg.ID("b")
			}).Op("=").ID("c")
		}),
		expect: `{
		a, b = c
		}`,
	},
	{
		desc: `list group`,
		code: BlockFunc(func(g *Group) { g.List(ID("a"), ID("b")).Op("=").ID("c") }),
		expect: `{
		a, b = c
		}`,
	},
	{
		desc:   `parens func`,
		code:   Parens(Lit(1)),
		expect: `(1)`,
	},
	{
		desc: `parens group`,
		code: BlockFunc(func(g *Group) { g.Parens(Lit(1)) }),
		expect: `{
		(1)
		}`,
	},
}

func TestGen(t *testing.T) {
	caseTester(t, gencases)
}
