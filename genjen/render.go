package main

import (
	"io"
	"strings"

	. "github.com/veggiemonk/jennifer/jen"
)

func render(w io.Writer) error {
	file := NewFile("jen")

	file.HeaderComment("This file is generated - do not edit.")
	file.Line()

	for _, b := range groups {
		comment := Commentf("%s %s", b.name, b.comment)

		if b.variadic && len(b.parameters) > 1 {
			panic("should not have variadic function with multiple params")
		}

		var variadic Code
		if b.variadic {
			variadic = Op("...")
		}
		var funcParams []Code
		var callParams []Code
		for _, name := range b.parameters {
			funcParams = append(funcParams, ID(name).Add(variadic).ID("Code"))
			callParams = append(callParams, ID(name).Add(variadic))
		}

		addFunctionAndGroupMethod(
			file,
			b.name,
			comment,
			funcParams,
			callParams,
			false,
		)

		/*
			// <comment>
			func (s *Statement) <name>(<funcParams>) *Statement {
				g := &Group{
					items:     []Code{<paramNames>}|<paramNames[0]>,
					name:      "<name>",
					open:      "<opening>",
					close:     "<closing>",
					separator: "<separator>",
					multi:     <multi>,
				}
				*s = append(*s, g)
				return s
			}
		*/
		file.Add(comment)
		file.Func().Params(
			ID("s").Op("*").ID("Statement"),
		).ID(b.name).Params(
			funcParams...,
		).Op("*").ID("Statement").Block(
			ID("g").Op(":=").Op("&").ID("Group").Values(Dict{
				ID("items"): Do(func(s *Statement) {
					if b.variadic {
						s.ID(b.parameters[0])
					} else {
						s.Index().ID("Code").ValuesFunc(func(g *Group) {
							for _, name := range b.parameters {
								g.ID(name)
							}
						})
					}
				}),
				ID("name"):      Lit(strings.ToLower(b.name)),
				ID("open"):      Lit(b.opening),
				ID("close"):     Lit(b.closing),
				ID("separator"): Lit(b.separator),
				ID("multi"):     Lit(b.multi),
			}),
			Op("*").ID("s").Op("=").Append(Op("*").ID("s"), ID("g")),
			Return(ID("s")),
		)

		if b.variadic && !b.preventFunc {

			funcName := b.name + "Func"
			funcComment := Commentf("%sFunc %s", b.name, b.comment)
			funcFuncParams := []Code{ID("f").Func().Params(Op("*").ID("Group"))}
			funcCallParams := []Code{ID("f")}

			addFunctionAndGroupMethod(
				file,
				funcName,
				funcComment,
				funcFuncParams,
				funcCallParams,
				false,
			)

			/*
				// <funcComment>
				func (s *Statement) <funcName>(f func(*Group)) *Statement {
					g := &Group{
						name:      "<name>",
						open:      "<opening>",
						close:     "<closing>",
						separator: "<separator>",
						multi:     <multi>,
					}
					f(g)
					*s = append(*s, g)
					return s
				}
			*/
			file.Add(funcComment)
			file.Func().Params(
				ID("s").Op("*").ID("Statement"),
			).ID(funcName).Params(
				funcFuncParams...,
			).Op("*").ID("Statement").Block(
				ID("g").Op(":=").Op("&").ID("Group").Values(Dict{
					ID("name"):      Lit(strings.ToLower(b.name)),
					ID("open"):      Lit(b.opening),
					ID("close"):     Lit(b.closing),
					ID("separator"): Lit(b.separator),
					ID("multi"):     Lit(b.multi),
				}),
				ID("f").Call(ID("g")),
				Op("*").ID("s").Op("=").Append(Op("*").ID("s"), ID("g")),
				Return(ID("s")),
			)
		}
	}

	type tkn struct {
		token     string
		name      string
		tokenType string
		tokenDesc string
	}
	tokens := []tkn{}
	for _, v := range identifiers {
		tokens = append(tokens, tkn{
			token:     v,
			name:      strings.ToUpper(v[:1]) + v[1:],
			tokenType: "identifierToken",
			tokenDesc: "identifier",
		})
	}
	for _, v := range keywords {
		tokens = append(tokens, tkn{
			token:     v,
			name:      strings.ToUpper(v[:1]) + v[1:],
			tokenType: "keywordToken",
			tokenDesc: "keyword",
		})
	}

	for i, t := range tokens {
		comment := Commentf(
			"%s renders the %s %s.",
			t.name,
			t.token,
			t.tokenDesc,
		)
		addFunctionAndGroupMethod(
			file,
			t.name,
			comment,
			nil,
			nil,
			i != 0, // only enforce test coverage on one item
		)

		/*
			// <comment>
			func (s *Statement) <name>() *Statement {
				t := token{
					typ:     <tokenType>,
					content: "<token>",
				}
				*s = append(*s, t)
				return s
			}
		*/
		file.Add(comment)
		file.Func().Params(
			ID("s").Op("*").ID("Statement"),
		).ID(t.name).Params().Op("*").ID("Statement").Block(
			Do(func(s *Statement) {
				if i != 0 {
					// only enforce test coverage on one item
					s.Comment("notest")
				}
			}),
			ID("t").Op(":=").ID("token").Values(Dict{
				ID("typ"):     ID(t.tokenType),
				ID("content"): Lit(t.token),
			}),
			Op("*").ID("s").Op("=").Append(Op("*").ID("s"), ID("t")),
			Return(ID("s")),
		)
	}

	return file.Render(w)
}

// For each method on *Statement, this generates a package level
// function and a method on *Group, both with the same name.
func addFunctionAndGroupMethod(
	file *File,
	name string,
	comment *Statement,
	funcParams []Code,
	callParams []Code,
	notest bool,
) {
	/*
		// <comment>
		func <name>(<funcParams>) *Statement {
			return newStatement().<name>(<callParams>)
		}
	*/
	file.Add(comment)
	file.Func().ID(name).Params(funcParams...).Op("*").ID("Statement").Block(
		Do(func(s *Statement) {
			if notest {
				// only enforce test coverage on one item
				s.Comment("notest")
			}
		}),
		Return(ID("newStatement").Call().Dot(name).Call(callParams...)),
	)
	/*
		// <comment>
		func (g *Group) <name>(<funcParams>) *Statement {
			s := <name>(<callParams>)
			g.items = append(g.items, s)
			return s
		}
	*/
	file.Add(comment)
	file.Func().Params(
		ID("g").Op("*").ID("Group"),
	).ID(name).Params(funcParams...).Op("*").ID("Statement").Block(
		Do(func(s *Statement) {
			if notest {
				// only enforce test coverage on one item
				s.Comment("notest")
			}
		}),
		ID("s").Op(":=").ID(name).Params(callParams...),
		ID("g").Dot("items").Op("=").Append(ID("g").Dot("items"), ID("s")),
		Return(ID("s")),
	)
}
