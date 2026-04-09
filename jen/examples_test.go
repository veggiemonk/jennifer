package jen_test

import (
	"bytes"
	"fmt"

	. "github.com/veggiemonk/jennifer/jen"
)

func ExampleTypesFunc_empty() {
	c := Func().ID("F").TypesFunc(func(group *Group) {}).Params().Block()
	fmt.Printf("%#v", c)
	// Output:
	// func F() {}
}

func ExampleTypesFunc_null() {
	c := Func().ID("F").TypesFunc(func(group *Group) {
		group.Null()
	}).Params().Block()
	fmt.Printf("%#v", c)
	// Output:
	// func F() {}
}

func ExampleTypes_empty() {
	c := Func().ID("F").Types().Params().Block()
	fmt.Printf("%#v", c)
	// Output:
	// func F() {}
}

func ExampleTypes_null() {
	c := Func().ID("F").Types(Null()).Params().Block()
	fmt.Printf("%#v", c)
	// Output:
	// func F() {}
}

func ExampleTypes_definition() {
	c := Func().ID("Keys").Types(
		ID("K").Comparable(),
		ID("V").Any(),
	).Params(
		ID("m").Map(ID("K")).ID("V"),
	).Index().ID("K").Block()
	fmt.Printf("%#v", c)
	// Output:
	// func Keys[K comparable, V any](m map[K]V) []K {}
}

func ExampleTypes_usage() {
	c := Return(ID("Keys").Types(Int(), String()).Call(ID("m")))
	fmt.Printf("%#v", c)
	// Output:
	// return Keys[int, string](m)
}

func ExampleUnion() {
	c := Type().ID("PredeclaredSignedInteger").Interface(
		Union(Int(), Int8(), Int16(), Int32(), Int64()),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type PredeclaredSignedInteger interface {
	//	int | int8 | int16 | int32 | int64
	// }
}

func ExampleOp_approximate() {
	c := Type().ID("AnyString").Interface(
		Op("~").String(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type AnyString interface {
	//	~string
	// }
}

func ExampleCase_blockWithMultipleStatements() {
	c := Switch(ID("a")).Block(
		Case(Lit(1)).Block(
			Var().ID("i").Int(),
			Var().ID("j").Int(),
		),
	)
	fmt.Printf("%#v", c)
	// Output:
	// switch a {
	// case 1:
	// 	var i int
	// 	var j int
	// }
}

func ExampleCustom() {
	multiLineCall := Options{
		Close:     ")",
		Multi:     true,
		Open:      "(",
		Separator: ",",
	}
	c := ID("foo").Custom(multiLineCall, Lit("a"), Lit("b"), Lit("c"))
	fmt.Printf("%#v", c)
	// Output:
	// foo(
	// 	"a",
	// 	"b",
	// 	"c",
	// )
}

func ExampleCustomFunc() {
	multiLineCall := Options{
		Close:     ")",
		Multi:     true,
		Open:      "(",
		Separator: ",",
	}
	c := ID("foo").CustomFunc(multiLineCall, func(g *Group) {
		g.Lit("a")
		g.Lit("b")
		g.Lit("c")
	})
	fmt.Printf("%#v", c)
	// Output:
	// foo(
	// 	"a",
	// 	"b",
	// 	"c",
	// )
}

func ExampleFile_ImportName_conflict() {
	f := NewFile("main")

	// We provide a hint that package foo/a should use name "a", but because package bar/a already
	// registers the required name, foo/a is aliased.
	f.ImportName("github.com/foo/a", "a")

	f.Func().ID("main").Params().Block(
		Qual("github.com/bar/a", "Bar").Call(),
		Qual("github.com/foo/a", "Foo").Call(),
	)
	fmt.Printf("%#v", f)

	// Output:
	// package main
	//
	// import (
	// 	a "github.com/bar/a"
	// 	a1 "github.com/foo/a"
	// )
	//
	// func main() {
	// 	a.Bar()
	// 	a1.Foo()
	// }
}

func ExampleFile_ImportAlias_conflict() {
	f := NewFile("main")

	// We provide a hint that package foo/a should use alias "b", but because package bar/b already
	// registers the required name, foo/a is aliased using the requested alias as a base.
	f.ImportName("github.com/foo/a", "b")

	f.Func().ID("main").Params().Block(
		Qual("github.com/bar/b", "Bar").Call(),
		Qual("github.com/foo/a", "Foo").Call(),
	)
	fmt.Printf("%#v", f)

	// Output:
	// package main
	//
	// import (
	// 	b "github.com/bar/b"
	// 	b1 "github.com/foo/a"
	// )
	//
	// func main() {
	// 	b.Bar()
	// 	b1.Foo()
	// }
}

func ExampleFile_ImportName() {
	f := NewFile("main")

	// package a should use name "a"
	f.ImportName("github.com/foo/a", "a")

	// package b is not used in the code so will not be included
	f.ImportName("github.com/foo/b", "b")

	f.Func().ID("main").Params().Block(
		Qual("github.com/foo/a", "A").Call(),
	)
	fmt.Printf("%#v", f)

	// Output:
	// package main
	//
	// import "github.com/foo/a"
	//
	// func main() {
	// 	a.A()
	// }
}

func ExampleFile_ImportNames() {
	// package a should use name "a", package b is not used in the code so will not be included
	names := map[string]string{
		"github.com/foo/a": "a",
		"github.com/foo/b": "b",
	}

	f := NewFile("main")
	f.ImportNames(names)
	f.Func().ID("main").Params().Block(
		Qual("github.com/foo/a", "A").Call(),
	)
	fmt.Printf("%#v", f)

	// Output:
	// package main
	//
	// import "github.com/foo/a"
	//
	// func main() {
	// 	a.A()
	// }
}

func ExampleFile_ImportAlias() {
	f := NewFile("main")

	// package a should be aliased to "b"
	f.ImportAlias("github.com/foo/a", "b")

	// package c is not used in the code so will not be included
	f.ImportAlias("github.com/foo/c", "c")

	f.Func().ID("main").Params().Block(
		Qual("github.com/foo/a", "A").Call(),
	)
	fmt.Printf("%#v", f)

	// Output:
	// package main
	//
	// import b "github.com/foo/a"
	//
	// func main() {
	// 	b.A()
	// }
}

func ExampleFile_ImportAlias_dot() {
	f := NewFile("main")

	// package a should be a dot-import
	f.ImportAlias("github.com/foo/a", ".")

	// package b should be a dot-import
	f.ImportAlias("github.com/foo/b", ".")

	// package c is not used in the code so will not be included
	f.ImportAlias("github.com/foo/c", ".")

	f.Func().ID("main").Params().Block(
		Qual("github.com/foo/a", "A").Call(),
		Qual("github.com/foo/b", "B").Call(),
	)
	fmt.Printf("%#v", f)

	// Output:
	// package main
	//
	// import (
	// 	. "github.com/foo/a"
	// 	. "github.com/foo/b"
	// )
	//
	// func main() {
	// 	A()
	// 	B()
	// }
}

func ExampleFile_CgoPreamble() {
	f := NewFile("a")
	f.CgoPreamble(`#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
	printf("%s\n", s);
}
`)
	f.Func().ID("init").Params().Block(
		ID("cs").Op(":=").Qual("C", "CString").Call(Lit("Hello from stdio\n")),
		Qual("C", "myprint").Call(ID("cs")),
		Qual("C", "free").Call(Qual("unsafe", "Pointer").Parens(ID("cs"))),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import "unsafe"
	//
	// /*
	// #include <stdio.h>
	// #include <stdlib.h>
	//
	// void myprint(char* s) {
	// 	printf("%s\n", s);
	// }
	// */
	// import "C"
	//
	// func init() {
	// 	cs := C.CString("Hello from stdio\n")
	// 	C.myprint(cs)
	// 	C.free(unsafe.Pointer(cs))
	// }
}

func ExampleFile_CgoPreamble_anon() {
	f := NewFile("a")
	f.CgoPreamble(`#include <stdio.h>`)
	f.Func().ID("init").Params().Block(
		Qual("foo.bar/a", "A"),
		Qual("foo.bar/b", "B"),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import (
	// 	a "foo.bar/a"
	// 	b "foo.bar/b"
	// )
	//
	// // #include <stdio.h>
	// import "C"
	//
	// func init() {
	// 	a.A
	// 	b.B
	// }
}

func ExampleFile_CgoPreamble_no_preamble() {
	f := NewFile("a")
	f.Func().ID("init").Params().Block(
		Qual("C", "Foo").Call(),
		Qual("fmt", "Print").Call(),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import (
	// 	"C"
	// 	"fmt"
	// )
	//
	// func init() {
	// 	C.Foo()
	// 	fmt.Print()
	// }
}

func ExampleFile_CgoPreamble_no_preamble_single() {
	f := NewFile("a")
	f.Func().ID("init").Params().Block(
		Qual("C", "Foo").Call(),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import "C"
	//
	// func init() {
	// 	C.Foo()
	// }
}

func ExampleFile_CgoPreamble_no_preamble_single_anon() {
	f := NewFile("a")
	f.Anon("C")
	f.Func().ID("init").Params().Block()
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import "C"
	//
	// func init() {}
}

func ExampleFile_CgoPreamble_no_preamble_anon() {
	f := NewFile("a")
	f.Anon("C")
	f.Func().ID("init").Params().Block(
		Qual("fmt", "Print").Call(),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import (
	// 	"C"
	// 	"fmt"
	// )
	//
	// func init() {
	// 	fmt.Print()
	// }
}

func ExampleOp_complex_conditions() {
	c := If(Parens(ID("a").Op("||").ID("b")).Op("&&").ID("c")).Block()
	fmt.Printf("%#v", c)
	// Output:
	// if (a || b) && c {
	// }
}

func ExampleLit_bool_true() {
	c := Lit(true)
	fmt.Printf("%#v", c)
	// Output:
	// true
}

func ExampleLit_bool_false() {
	c := Lit(false)
	fmt.Printf("%#v", c)
	// Output:
	// false
}

func ExampleLit_byte() {
	// Lit can't tell the difference between byte and uint8. Use LitByte to
	// render byte literals.
	c := Lit(byte(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uint8(0x1)
}

func ExampleLit_complex64() {
	c := Lit(complex64(0 + 0i))
	fmt.Printf("%#v", c)
	// Output:
	// complex64(0 + 0i)
}

func ExampleLit_complex128() {
	c := Lit(0 + 0i)
	fmt.Printf("%#v", c)
	// Output:
	// (0 + 0i)
}

func ExampleLit_float32() {
	c := Lit(float32(1))
	fmt.Printf("%#v", c)
	// Output:
	// float32(1)
}

func ExampleLit_float64_one_point_zero() {
	c := Lit(1.0)
	fmt.Printf("%#v", c)
	// Output:
	// 1.0
}

func ExampleLit_float64_zero() {
	c := Lit(0.0)
	fmt.Printf("%#v", c)
	// Output:
	// 0.0
}

func ExampleLit_float64_negative() {
	c := Lit(-0.1)
	fmt.Printf("%#v", c)
	// Output:
	// -0.1
}

func ExampleLit_float64_negative_whole() {
	c := Lit(-1.0)
	fmt.Printf("%#v", c)
	// Output:
	// -1.0
}

func ExampleLit_int() {
	c := Lit(1)
	fmt.Printf("%#v", c)
	// Output:
	// 1
}

func ExampleLit_int8() {
	c := Lit(int8(1))
	fmt.Printf("%#v", c)
	// Output:
	// int8(1)
}

func ExampleLit_int16() {
	c := Lit(int16(1))
	fmt.Printf("%#v", c)
	// Output:
	// int16(1)
}

func ExampleLit_int32() {
	c := Lit(int32(1))
	fmt.Printf("%#v", c)
	// Output:
	// int32(1)
}

func ExampleLit_int64() {
	c := Lit(int64(1))
	fmt.Printf("%#v", c)
	// Output:
	// int64(1)
}

func ExampleLit_uint() {
	c := Lit(uint(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uint(0x1)
}

func ExampleLit_uint8() {
	c := Lit(uint8(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uint8(0x1)
}

func ExampleLit_uint16() {
	c := Lit(uint16(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uint16(0x1)
}

func ExampleLit_uint32() {
	c := Lit(uint32(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uint32(0x1)
}

func ExampleLit_uint64() {
	c := Lit(uint64(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uint64(0x1)
}

func ExampleLit_uintptr() {
	c := Lit(uintptr(0x1))
	fmt.Printf("%#v", c)
	// Output:
	// uintptr(0x1)
}

func ExampleLit_rune() {
	// Lit can't tell the difference between rune and int32. Use LitRune to
	// render rune literals.
	c := Lit('x')
	fmt.Printf("%#v", c)
	// Output:
	// int32(120)
}

func ExampleLitRune() {
	c := LitRune('x')
	fmt.Printf("%#v", c)
	// Output:
	// 'x'
}

func ExampleLitRuneFunc() {
	c := LitRuneFunc(func() rune {
		return '\t'
	})
	fmt.Printf("%#v", c)
	// Output:
	// '\t'
}

func ExampleLitByte() {
	c := LitByte(byte(1))
	fmt.Printf("%#v", c)
	// Output:
	// byte(0x1)
}

func ExampleLitByteFunc() {
	c := LitByteFunc(func() byte {
		return byte(2)
	})
	fmt.Printf("%#v", c)
	// Output:
	// byte(0x2)
}

func ExampleLit_string() {
	c := Lit("foo")
	fmt.Printf("%#v", c)
	// Output:
	// "foo"
}

func ExampleValues_dict_single() {
	c := Map(String()).String().Values(Dict{
		Lit("a"): Lit("b"),
	})
	fmt.Printf("%#v", c)
	// Output:
	// map[string]string{"a": "b"}
}

func ExampleValues_dict_multiple() {
	c := Map(String()).String().Values(Dict{
		Lit("a"): Lit("b"),
		Lit("c"): Lit("d"),
	})
	fmt.Printf("%#v", c)
	// Output:
	// map[string]string{
	// 	"a": "b",
	// 	"c": "d",
	// }
}

func ExampleValues_dict_composite() {
	c := Op("&").ID("Person").Values(Dict{
		ID("Age"):  Lit(1),
		ID("Name"): Lit("a"),
	})
	fmt.Printf("%#v", c)
	// Output:
	// &Person{
	// 	Age:  1,
	// 	Name: "a",
	// }
}

func ExampleAdd() {
	ptr := Op("*")
	c := ID("a").Op("=").Add(ptr).ID("b")
	fmt.Printf("%#v", c)
	// Output:
	// a = *b
}

func ExampleAdd_var() {
	a := ID("a")
	i := Int()
	c := Var().Add(a, i)
	fmt.Printf("%#v", c)
	// Output:
	// var a int
}

func ExampleAppend() {
	c := Append(ID("a"), ID("b"))
	fmt.Printf("%#v", c)
	// Output:
	// append(a, b)
}

func ExampleAppend_more() {
	c := ID("a").Op("=").Append(ID("a"), ID("b").Op("..."))
	fmt.Printf("%#v", c)
	// Output:
	// a = append(a, b...)
}

func ExampleAssert() {
	c := List(ID("b"), ID("ok")).Op(":=").ID("a").Assert(Bool())
	fmt.Printf("%#v", c)
	// Output:
	// b, ok := a.(bool)
}

func ExampleBlock() {
	c := Func().ID("foo").Params().String().Block(
		ID("a").Op("=").ID("b"),
		ID("b").Op("++"),
		Return(ID("b")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func foo() string {
	// 	a = b
	// 	b++
	// 	return b
	// }
}

func ExampleBlock_if() {
	c := If(ID("a").Op(">").Lit(10)).Block(
		ID("a").Op("=").ID("a").Op("/").Lit(2),
	)
	fmt.Printf("%#v", c)
	// Output:
	// if a > 10 {
	// 	a = a / 2
	// }
}

func ExampleValuesFunc() {
	c := ID("numbers").Op(":=").Index().Int().ValuesFunc(func(g *Group) {
		for i := 0; i <= 5; i++ {
			g.Lit(i)
		}
	})
	fmt.Printf("%#v", c)
	// Output:
	// numbers := []int{0, 1, 2, 3, 4, 5}
}

func ExampleBlockFunc() {
	increment := true
	name := "a"
	c := Func().ID("a").Params().BlockFunc(func(g *Group) {
		g.ID(name).Op("=").Lit(1)
		if increment {
			g.ID(name).Op("++")
		} else {
			g.ID(name).Op("--")
		}
	})
	fmt.Printf("%#v", c)
	// Output:
	// func a() {
	// 	a = 1
	// 	a++
	// }
}

func ExampleBool() {
	c := Var().ID("b").Bool()
	fmt.Printf("%#v", c)
	// Output:
	// var b bool
}

func ExampleBreak() {
	c := For(
		ID("i").Op(":=").Lit(0),
		ID("i").Op("<").Lit(10),
		ID("i").Op("++"),
	).Block(
		If(ID("i").Op(">").Lit(5)).Block(
			Break(),
		),
	)
	fmt.Printf("%#v", c)
	// Output:
	// for i := 0; i < 10; i++ {
	// 	if i > 5 {
	// 		break
	// 	}
	// }
}

func ExampleByte() {
	c := ID("b").Op(":=").ID("a").Assert(Byte())
	fmt.Printf("%#v", c)
	// Output:
	// b := a.(byte)
}

func ExampleCall() {
	c := Qual("fmt", "Printf").Call(
		Lit("%#v: %T\n"),
		ID("a"),
		ID("b"),
	)
	fmt.Printf("%#v", c)
	// Output:
	// fmt.Printf("%#v: %T\n", a, b)
}

func ExampleCall_fmt() {
	c := ID("a").Call(Lit("b"))
	fmt.Printf("%#v", c)
	// Output:
	// a("b")
}

func ExampleCallFunc() {
	f := func(name, second string) {
		c := ID("foo").CallFunc(func(g *Group) {
			g.ID(name)
			if second != "" {
				g.Lit(second)
			}
		})
		fmt.Printf("%#v\n", c)
	}
	f("a", "b")
	f("c", "")
	// Output:
	// foo(a, "b")
	// foo(c)
}

func ExampleCap() {
	c := ID("i").Op(":=").Cap(ID("v"))
	fmt.Printf("%#v", c)
	// Output:
	// i := cap(v)
}

func ExampleCase() {
	c := Switch(ID("person")).Block(
		Case(ID("John"), ID("Peter")).Block(
			Return(Lit("male")),
		),
		Case(ID("Gill")).Block(
			Return(Lit("female")),
		),
	)
	fmt.Printf("%#v", c)
	// Output:
	// switch person {
	// case John, Peter:
	// 	return "male"
	// case Gill:
	// 	return "female"
	// }
}

func ExampleBlock_case() {
	c := Select().Block(
		Case(Op("<-").ID("done")).Block(
			Return(Nil()),
		),
		Case(List(Err(), ID("open")).Op(":=").Op("<-").ID("fail")).Block(
			If(Op("!").ID("open")).Block(
				Return(Err()),
			),
		),
	)
	fmt.Printf("%#v", c)
	// Output:
	// select {
	// case <-done:
	// 	return nil
	// case err, open := <-fail:
	// 	if !open {
	// 		return err
	// 	}
	// }
}

func ExampleBlockFunc_case() {
	preventExitOnError := true
	c := Select().Block(
		Case(Op("<-").ID("done")).Block(
			Return(Nil()),
		),
		Case(Err().Op(":=").Op("<-").ID("fail")).BlockFunc(func(g *Group) {
			if !preventExitOnError {
				g.Return(Err())
			} else {
				g.Qual("fmt", "Println").Call(Err())
			}
		}),
	)
	fmt.Printf("%#v", c)
	// Output:
	// select {
	// case <-done:
	// 	return nil
	// case err := <-fail:
	// 	fmt.Println(err)
	// }
}

func ExampleCaseFunc() {
	samIsMale := false
	c := Switch(ID("person")).Block(
		CaseFunc(func(g *Group) {
			g.ID("John")
			g.ID("Peter")
			if samIsMale {
				g.ID("Sam")
			}
		}).Block(
			Return(Lit("male")),
		),
		CaseFunc(func(g *Group) {
			g.ID("Gill")
			if !samIsMale {
				g.ID("Sam")
			}
		}).Block(
			Return(Lit("female")),
		),
	)
	fmt.Printf("%#v", c)
	// Output:
	// switch person {
	// case John, Peter:
	// 	return "male"
	// case Gill, Sam:
	// 	return "female"
	// }
}

func ExampleChan() {
	c := Func().ID("init").Params().Block(
		ID("c").Op(":=").Make(Chan().Qual("os", "Signal"), Lit(1)),
		Qual("os/signal", "Notify").Call(ID("c"), Qual("os", "Interrupt")),
		Qual("os/signal", "Notify").Call(ID("c"), Qual("syscall", "SIGTERM")),
		Go().Func().Params().Block(
			Op("<-").ID("c"),
			ID("cancel").Call(),
		).Call(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func init() {
	// 	c := make(chan os.Signal, 1)
	// 	signal.Notify(c, os.Interrupt)
	// 	signal.Notify(c, syscall.SIGTERM)
	// 	go func() {
	// 		<-c
	// 		cancel()
	// 	}()
	// }
}

func ExampleClose() {
	c := Block(
		ID("ch").Op(":=").Make(Chan().Struct()),
		Go().Func().Params().Block(
			Op("<-").ID("ch"),
			Qual("fmt", "Println").Call(Lit("done.")),
		).Call(),
		Close(ID("ch")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// {
	// 	ch := make(chan struct{})
	// 	go func() {
	// 		<-ch
	// 		fmt.Println("done.")
	// 	}()
	// 	close(ch)
	// }
}

func ExampleClear() {
	c := Block(
		ID("a").Op(":=").Map(String()).String().Values(),
		Clear(ID("a")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// {
	// 	a := map[string]string{}
	// 	clear(a)
	// }
}

func ExampleMin() {
	c := Block(
		ID("n").Op(":=").Min(Lit(1), Lit(2)),
	)
	fmt.Printf("%#v", c)
	// Output:
	// {
	// 	n := min(1, 2)
	// }
}

func ExampleMax() {
	c := Block(
		ID("x").Op(":=").Max(Lit(1), Lit(2)),
	)
	fmt.Printf("%#v", c)
	// Output:
	// {
	// 	x := max(1, 2)
	// }
}

func ExampleComment() {
	f := NewFile("a")
	f.Comment("Foo returns the string \"foo\"")
	f.Func().ID("Foo").Params().String().Block(
		Return(Lit("foo")).Comment("return the string foo"),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// // Foo returns the string "foo"
	// func Foo() string {
	// 	return "foo" // return the string foo
	// }
}

func ExampleComment_multiline() {
	c := Comment("a\nb")
	fmt.Printf("%#v", c)
	// Output:
	// /*
	// a
	// b
	// */
}

func ExampleComment_formatting_disabled() {
	c := ID("foo").Call(Comment("/* inline */")).Comment("//no-space")
	fmt.Printf("%#v", c)
	// Output:
	// foo( /* inline */ ) //no-space
}

func ExampleCommentf() {
	name := "foo"
	val := "bar"
	c := ID(name).Op(":=").Lit(val).Commentf("%s is the string \"%s\"", name, val)
	fmt.Printf("%#v", c)
	// Output:
	// foo := "bar" // foo is the string "bar"
}

func ExampleComplex() {
	c := Func().ID("main").Params().Block(
		ID("c1").Op(":=").Lit(1+3.75i),
		ID("c2").Op(":=").Complex(Lit(1), Lit(3.75)),
		Qual("fmt", "Println").Call(ID("c1").Op("==").ID("c2")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func main() {
	// 	c1 := (1 + 3.75i)
	// 	c2 := complex(1, 3.75)
	// 	fmt.Println(c1 == c2)
	// }
}

func ExampleComplex128() {
	c := Func().ID("main").Params().Block(
		Var().ID("c").Complex128(),
		ID("c").Op("=").Lit(1+2i),
		Qual("fmt", "Println").Call(ID("c")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func main() {
	// 	var c complex128
	// 	c = (1 + 2i)
	// 	fmt.Println(c)
	// }
}

func ExampleComplex64() {
	c := Func().ID("main").Params().Block(
		Var().ID("c64").Complex64(),
		ID("c64").Op("=").Complex(Lit(5), Float32().Parens(Lit(2))),
		Qual("fmt", "Printf").Call(Lit("%T\n"), ID("c64")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func main() {
	// 	var c64 complex64
	// 	c64 = complex(5, float32(2))
	// 	fmt.Printf("%T\n", c64)
	// }
}

func ExampleParams() {
	c := Func().Params(
		ID("a").ID("A"),
	).ID("foo").Params(
		ID("b"),
		ID("c").String(),
	).String().Block(
		Return(ID("b").Op("+").ID("c")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func (a A) foo(b, c string) string {
	// 	return b + c
	// }
}

func ExampleIndex() {
	c := Var().ID("a").Index().String()
	fmt.Printf("%#v", c)
	// Output:
	// var a []string
}

func ExampleIndex_index() {
	c := ID("a").Op(":=").ID("b").Index(Lit(0), Lit(1))
	fmt.Printf("%#v", c)
	// Output:
	// a := b[0:1]
}

func ExampleIndex_empty() {
	c := ID("a").Op(":=").ID("b").Index(Lit(1), Empty())
	fmt.Printf("%#v", c)
	// Output:
	// a := b[1:]
}

func ExampleOp() {
	c := ID("a").Op(":=").ID("b").Call()
	fmt.Printf("%#v", c)
	// Output:
	// a := b()
}

func ExampleOp_star() {
	c := ID("a").Op("=").Op("*").ID("b")
	fmt.Printf("%#v", c)
	// Output:
	// a = *b
}

func ExampleOp_variadic() {
	c := ID("a").Call(ID("b").Op("..."))
	fmt.Printf("%#v", c)
	// Output:
	// a(b...)
}

func ExampleNewFilePath() {
	f := NewFilePath("a.b/c")
	f.Func().ID("init").Params().Block(
		Qual("a.b/c", "Foo").Call().Comment("Local package - alias is omitted."),
		Qual("d.e/f", "Bar").Call().Comment("Import is automatically added."),
		Qual("g.h/f", "Baz").Call().Comment("Colliding package name is automatically renamed."),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package c
	//
	// import (
	// 	f "d.e/f"
	// 	f1 "g.h/f"
	// )
	//
	// func init() {
	// 	Foo()    // Local package - alias is omitted.
	// 	f.Bar()  // Import is automatically added.
	// 	f1.Baz() // Colliding package name is automatically renamed.
	// }
}

func ExampleStruct_empty() {
	c := ID("c").Op(":=").Make(Chan().Struct())
	fmt.Printf("%#v", c)
	// Output:
	// c := make(chan struct{})
}

func ExampleStruct() {
	c := Type().ID("foo").Struct(
		List(ID("x"), ID("y")).Int(),
		ID("u").Float32(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type foo struct {
	// 	x, y int
	// 	u    float32
	// }
}

func ExampleDefer() {
	c := Defer().ID("foo").Call()
	fmt.Printf("%#v", c)
	// Output:
	// defer foo()
}

func ExampleGoto() {
	c := Goto().ID("Outer")
	fmt.Printf("%#v", c)
	// Output:
	// goto Outer
}

func ExampleStatement_Clone_broken() {
	a := ID("a")
	c := Block(
		a.Call(),
		a.Call(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// {
	// 	a()()
	// 	a()()
	// }
}

func ExampleStatement_Clone_fixed() {
	a := ID("a")
	c := Block(
		a.Clone().Call(),
		a.Clone().Call(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// {
	// 	a()
	// 	a()
	// }
}

func ExampleFile_Render() {
	f := NewFile("a")
	f.Func().ID("main").Params().Block()
	buf := &bytes.Buffer{}
	err := f.Render(buf)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(buf.String())
	}
	// Output:
	// package a
	//
	// func main() {}
}

func ExampleLit() {
	c := ID("a").Op(":=").Lit("a")
	fmt.Printf("%#v", c)
	// Output:
	// a := "a"
}

func ExampleLit_float() {
	c := ID("a").Op(":=").Lit(1.5)
	fmt.Printf("%#v", c)
	// Output:
	// a := 1.5
}

func ExampleLitFunc() {
	c := ID("a").Op(":=").LitFunc(func() any { return 1 + 1 })
	fmt.Printf("%#v", c)
	// Output:
	// a := 2
}

func ExampleDot() {
	c := Qual("a.b/c", "Foo").Call().Dot("Bar").Index(Lit(0)).Dot("Baz")
	fmt.Printf("%#v", c)
	// Output:
	// c.Foo().Bar[0].Baz
}

func ExampleList() {
	c := List(ID("a"), Err()).Op(":=").ID("b").Call()
	fmt.Printf("%#v", c)
	// Output:
	// a, err := b()
}

func ExampleQual() {
	c := Qual("encoding/gob", "NewEncoder").Call()
	fmt.Printf("%#v", c)
	// Output:
	// gob.NewEncoder()
}

func ExampleQual_file() {
	f := NewFilePath("a.b/c")
	f.Func().ID("init").Params().Block(
		Qual("a.b/c", "Foo").Call().Comment("Local package - name is omitted."),
		Qual("d.e/f", "Bar").Call().Comment("Import is automatically added."),
		Qual("g.h/f", "Baz").Call().Comment("Colliding package name is renamed."),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package c
	//
	// import (
	// 	f "d.e/f"
	// 	f1 "g.h/f"
	// )
	//
	// func init() {
	// 	Foo()    // Local package - name is omitted.
	// 	f.Bar()  // Import is automatically added.
	// 	f1.Baz() // Colliding package name is renamed.
	// }
}

func ExampleQual_local() {
	f := NewFilePath("a.b/c")
	f.Func().ID("main").Params().Block(
		Qual("a.b/c", "D").Call(),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package c
	//
	// func main() {
	// 	D()
	// }
}

func ExampleID() {
	c := If(ID("i").Op("==").ID("j")).Block(
		Return(ID("i")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// if i == j {
	// 	return i
	// }
}

func ExampleErr() {
	c := If(
		Err().Op(":=").ID("foo").Call(),
		Err().Op("!=").Nil(),
	).Block(
		Return(Err()),
	)
	fmt.Printf("%#v", c)
	// Output:
	// if err := foo(); err != nil {
	// 	return err
	// }
}

func ExampleSwitch() {
	c := Switch(ID("value").Dot("Kind").Call()).Block(
		Case(Qual("reflect", "Float32"), Qual("reflect", "Float64")).Block(
			Return(Lit("float")),
		),
		Case(Qual("reflect", "Bool")).Block(
			Return(Lit("bool")),
		),
		Case(Qual("reflect", "Uintptr")).Block(
			Fallthrough(),
		),
		Default().Block(
			Return(Lit("none")),
		),
	)
	fmt.Printf("%#v", c)
	// Output:
	// switch value.Kind() {
	// case reflect.Float32, reflect.Float64:
	// 	return "float"
	// case reflect.Bool:
	// 	return "bool"
	// case reflect.Uintptr:
	// 	fallthrough
	// default:
	// 	return "none"
	// }
}

func ExampleTag() {
	c := Type().ID("foo").Struct(
		ID("A").String().Tag(map[string]string{"json": "a"}),
		ID("B").Int().Tag(map[string]string{"json": "b", "bar": "baz"}),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type foo struct {
	// 	A string `json:"a"`
	// 	B int    `bar:"baz" json:"b"`
	// }
}

func ExampleTag_withQuotesAndNewline() {
	c := Type().ID("foo").Struct(
		ID("A").String().Tag(map[string]string{"json": "a"}),
		ID("B").Int().Tag(map[string]string{"json": "b", "bar": "the value of\nthe\"bar\" tag"}),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type foo struct {
	// 	A string `json:"a"`
	// 	B int    `bar:"the value of\nthe\"bar\" tag" json:"b"`
	// }
}

func ExampleNull_and_nil() {
	c := Func().ID("foo").Params(
		nil,
		ID("s").String(),
		Null(),
		ID("i").Int(),
	).Block()
	fmt.Printf("%#v", c)
	// Output:
	// func foo(s string, i int) {}
}

func ExampleNull_index() {
	c := ID("a").Op(":=").ID("b").Index(Lit(1), Null())
	fmt.Printf("%#v", c)
	// Output:
	// a := b[1]
}

func ExampleEmpty() {
	c := ID("a").Op(":=").ID("b").Index(Lit(1), Empty())
	fmt.Printf("%#v", c)
	// Output:
	// a := b[1:]
}

func ExampleBlock_complex() {
	collection := func(name string, key, value Code) *Statement {
		if key == nil {
			// slice
			return Var().ID(name).Index().Add(value)
		} else {
			// map
			return Var().ID(name).Map(key).Add(value)
		}
	}
	c := Func().ID("main").Params().Block(
		collection("foo", nil, String()),
		collection("bar", String(), Int()),
	)
	fmt.Printf("%#v", c)
	// Output:
	// func main() {
	// 	var foo []string
	// 	var bar map[string]int
	// }
}

func ExampleFunc_declaration() {
	c := Func().ID("a").Params().Block()
	fmt.Printf("%#v", c)
	// Output:
	// func a() {}
}

func ExampleFunc_literal() {
	c := ID("a").Op(":=").Func().Params().Block()
	fmt.Printf("%#v", c)
	// Output:
	// a := func() {}
}

func ExampleInterface() {
	c := Type().ID("a").Interface(
		ID("b").Params().String(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type a interface {
	// 	b() string
	// }
}

func ExampleInterface_empty() {
	c := Var().ID("a").Interface()
	fmt.Printf("%#v", c)
	// Output:
	// var a interface{}
}

func ExampleParens() {
	c := ID("b").Op(":=").Index().Byte().Parens(ID("s"))
	fmt.Printf("%#v", c)
	// Output:
	// b := []byte(s)
}

func ExampleParens_order() {
	c := ID("a").Op("/").Parens(ID("b").Op("+").ID("c"))
	fmt.Printf("%#v", c)
	// Output:
	// a / (b + c)
}

func ExampleValues() {
	c := Index().String().Values(Lit("a"), Lit("b"))
	fmt.Printf("%#v", c)
	// Output:
	// []string{"a", "b"}
}

func ExampleDo() {
	f := func(name string, isMap bool) *Statement {
		return ID(name).Op(":=").Do(func(s *Statement) {
			if isMap {
				s.Map(String()).String()
			} else {
				s.Index().String()
			}
		}).Values()
	}
	fmt.Printf("%#v\n%#v", f("a", true), f("b", false))
	// Output:
	// a := map[string]string{}
	// b := []string{}
}

func ExampleReturn() {
	c := Return(ID("a"), ID("b"))
	fmt.Printf("%#v", c)
	// Output:
	// return a, b
}

func ExampleMap() {
	c := ID("a").Op(":=").Map(String()).String().Values()
	fmt.Printf("%#v", c)
	// Output:
	// a := map[string]string{}
}

func ExampleDict() {
	c := ID("a").Op(":=").Map(String()).String().Values(Dict{
		Lit("a"): Lit("b"),
		Lit("c"): Lit("d"),
	})
	fmt.Printf("%#v", c)
	// Output:
	// a := map[string]string{
	// 	"a": "b",
	// 	"c": "d",
	// }
}

func ExampleDict_nil() {
	c := ID("a").Op(":=").Map(String()).String().Values()
	fmt.Printf("%#v", c)
	// Output:
	// a := map[string]string{}
}

func ExampleDictFunc() {
	c := ID("a").Op(":=").Map(String()).String().Values(DictFunc(func(d Dict) {
		d[Lit("a")] = Lit("b")
		d[Lit("c")] = Lit("d")
	}))
	fmt.Printf("%#v", c)
	// Output:
	// a := map[string]string{
	// 	"a": "b",
	// 	"c": "d",
	// }
}

func ExampleDefs() {
	c := Const().Defs(
		ID("a").Op("=").Lit("a"),
		ID("b").Op("=").Lit("b"),
	)
	fmt.Printf("%#v", c)
	// Output:
	// const (
	// 	a = "a"
	// 	b = "b"
	// )
}

func ExampleIf() {
	c := If(
		Err().Op(":=").ID("a").Call(),
		Err().Op("!=").Nil(),
	).Block(
		Return(Err()),
	)
	fmt.Printf("%#v", c)
	// Output:
	// if err := a(); err != nil {
	// 	return err
	// }
}

func ExampleID_local() {
	c := ID("a").Op(":=").Lit(1)
	fmt.Printf("%#v", c)
	// Output:
	// a := 1
}

func ExampleID_select() {
	c := ID("a").Dot("b").Dot("c").Call()
	fmt.Printf("%#v", c)
	// Output:
	// a.b.c()
}

func ExampleID_remote() {
	f := NewFile("main")
	f.Func().ID("main").Params().Block(
		Qual("fmt", "Println").Call(
			Lit("Hello, world"),
		),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package main
	//
	// import "fmt"
	//
	// func main() {
	// 	fmt.Println("Hello, world")
	// }
}

func ExampleFor() {
	c := For(
		ID("i").Op(":=").Lit(0),
		ID("i").Op("<").Lit(10),
		ID("i").Op("++"),
	).Block(
		Qual("fmt", "Println").Call(ID("i")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(i)
	// }
}

func ExampleNewFile() {
	f := NewFile("main")
	f.Func().ID("main").Params().Block(
		Qual("fmt", "Println").Call(Lit("Hello, world")),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package main
	//
	// import "fmt"
	//
	// func main() {
	// 	fmt.Println("Hello, world")
	// }
}

func ExampleNewFilePathName() {
	f := NewFilePathName("a.b/c", "main")
	f.Func().ID("main").Params().Block(
		Qual("a.b/c", "Foo").Call(),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package main
	//
	// func main() {
	// 	Foo()
	// }
}

func ExampleFile_HeaderComment_withPackageComment() {
	f := NewFile("c")
	f.CanonicalPath = "d.e/f"
	f.HeaderComment("Code generated by...")
	f.PackageComment("Package c implements...")
	f.Func().ID("init").Params().Block()
	fmt.Printf("%#v", f)
	// Output:
	// // Code generated by...
	//
	// // Package c implements...
	// package c // import "d.e/f"
	//
	// func init() {}
}

func ExampleFile_Anon() {
	f := NewFile("c")
	f.Anon("a")
	f.Func().ID("init").Params().Block()
	fmt.Printf("%#v", f)
	// Output:
	// package c
	//
	// import _ "a"
	//
	// func init() {}
}

func ExampleFile_PackagePrefix() {
	f := NewFile("a")
	f.PackagePrefix = "pkg"
	f.Func().ID("main").Params().Block(
		Qual("b.c/d", "E").Call(),
	)
	fmt.Printf("%#v", f)
	// Output:
	// package a
	//
	// import pkg_d "b.c/d"
	//
	// func main() {
	// 	pkg_d.E()
	// }
}

func ExampleFile_NoFormat() {
	f := NewFile("main")
	f.NoFormat = true

	f.Func().ID("main").Params().Block(
		Qual("fmt", "Println").Call(Lit("foo")),
	)
	fmt.Printf("%q", fmt.Sprintf("%#v", f)) // Special case because Go Examples don't handle multiple newlines well.

	// Output:
	// "package main\n\nimport \"fmt\"\n\n\nfunc main () {\nfmt.Println (\"foo\")\n}"
}

// Go 1.22: range over integer
func ExampleFor_rangeOverInt() {
	c := For(
		ID("i").Op(":=").Range().Lit(10),
	).Block(
		Qual("fmt", "Println").Call(ID("i")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// for i := range 10 {
	// 	fmt.Println(i)
	// }
}

// Go 1.22: range over integer, no variable
func ExampleFor_rangeOverIntNoVar() {
	c := For(
		Range().Lit(10),
	).Block(
		Qual("fmt", "Println").Call(Lit("hello")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// for range 10 {
	// 	fmt.Println("hello")
	// }
}

// Go 1.23: range over iterator function (func(yield func(int) bool))
func ExampleFor_rangeOverFunc() {
	c := For(
		ID("v").Op(":=").Range().ID("seq"),
	).Block(
		Qual("fmt", "Println").Call(ID("v")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// for v := range seq {
	// 	fmt.Println(v)
	// }
}

// Go 1.23: range over two-value iterator function (func(yield func(int, string) bool))
func ExampleFor_rangeOverFunc2() {
	c := For(
		List(ID("k"), ID("v")).Op(":=").Range().ID("seq2"),
	).Block(
		Qual("fmt", "Println").Call(ID("k"), ID("v")),
	)
	fmt.Printf("%#v", c)
	// Output:
	// for k, v := range seq2 {
	// 	fmt.Println(k, v)
	// }
}

// Go 1.23: iterator function type signature
func ExampleFunc_iteratorType() {
	// type Seq = func(yield func(int) bool)
	c := Type().ID("Seq").Op("=").Func().Params(
		ID("yield").Func().Params(Int()).Bool(),
	)
	fmt.Printf("%#v", c)
	// Output:
	// type Seq = func(yield func(int) bool)
}

// Go 1.24: generic type alias
func ExampleType_genericAlias() {
	// type Set[T comparable] = map[T]struct{}
	c := Type().ID("Set").Types(ID("T").Comparable()).Op("=").Map(ID("T")).Struct()
	fmt.Printf("%#v", c)
	// Output:
	// type Set[T comparable] = map[T]struct{}
}
