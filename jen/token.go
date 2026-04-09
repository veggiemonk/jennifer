package jen

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type tokenType int

const (
	packageToken tokenType = iota
	identifierToken
	keywordToken
	operatorToken
	delimiterToken
	literalToken
	literalRuneToken
	literalByteToken
	nullToken
	layoutToken
)

type token struct {
	typ     tokenType
	content any
}

func (t token) isNull(f *File) bool {
	if t.typ == packageToken {
		return f.isDotImport(t.content.(string)) || f.isLocal(t.content.(string))
	}
	return t.typ == nullToken
}

func (t token) render(f *File, w io.Writer, s *Statement) error {
	switch t.typ {
	case literalToken:
		var out string
		switch t.content.(type) {
		case bool, string, int, complex128:
			out = fmt.Sprintf("%#v", t.content)
		case float64:
			out = fmt.Sprintf("%#v", t.content)
			if !strings.Contains(out, ".") && !strings.Contains(out, "e") {
				out += ".0"
			}
		case float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
			out = fmt.Sprintf("%T(%#v)", t.content, t.content)
		case complex64:
			out = fmt.Sprintf("%T%#v", t.content, t.content)
		default:
			panic(fmt.Sprintf("unsupported type for literal: %T", t.content))
		}
		if _, err := w.Write([]byte(out)); err != nil {
			return err
		}
	case literalRuneToken:
		if _, err := w.Write([]byte(strconv.QuoteRune(t.content.(rune)))); err != nil {
			return err
		}
	case literalByteToken:
		if _, err := fmt.Fprintf(w, "byte(%#v)", t.content); err != nil {
			return err
		}
	case keywordToken, operatorToken, layoutToken, delimiterToken:
		if _, err := fmt.Fprint(w, t.content); err != nil {
			return err
		}
	case packageToken:
		alias := f.register(t.content.(string))
		if _, err := w.Write([]byte(alias)); err != nil {
			return err
		}
	case identifierToken:
		if _, err := w.Write([]byte(t.content.(string))); err != nil {
			return err
		}
	case nullToken:
		// nothing
	}
	return nil
}

// qualToken renders a qualified identifier (package.Name) and registers the import.
type qualToken struct {
	path string
	name string
}

func (q qualToken) isNull(f *File) bool {
	return false
}

func (q qualToken) render(f *File, w io.Writer, s *Statement) error {
	if f.isDotImport(q.path) {
		_, err := w.Write([]byte(q.name))
		return err
	}
	if f.isLocal(q.path) {
		_, err := w.Write([]byte(q.name))
		return err
	}
	alias := f.register(q.path)
	_, err := fmt.Fprintf(w, "%s.%s", alias, q.name)
	return err
}

// --- Token-based builder methods ---

// newToken creates a new Statement with a single token.
func newToken(typ tokenType, content any) *Statement {
	s := newStatement()
	*s = append(*s, token{typ: typ, content: content})
	return s
}

// addToken appends a token to the statement.
func (s *Statement) addToken(typ tokenType, content any) *Statement {
	*s = append(*s, token{typ: typ, content: content})
	return s
}

// Null adds a null item. Null items render nothing and are not followed by a
// separator in lists.
func Null() *Statement                { return newToken(nullToken, nil) }
func (g *Group) Null() *Statement     { return g.item(Null()) }
func (s *Statement) Null() *Statement { return s.addToken(nullToken, nil) }

// Empty adds an empty item. Empty items render nothing but are followed by a
// separator in lists.
func Empty() *Statement                { return newToken(operatorToken, "") }
func (g *Group) Empty() *Statement     { return g.item(Empty()) }
func (s *Statement) Empty() *Statement { return s.addToken(operatorToken, "") }

// Op renders the provided operator / token.
func Op(op string) *Statement                { return newToken(operatorToken, op) }
func (g *Group) Op(op string) *Statement     { return g.item(Op(op)) }
func (s *Statement) Op(op string) *Statement { return s.addToken(operatorToken, op) }

// Dot renders a period followed by an identifier. Use for fields and selectors.
func Dot(name string) *Statement            { return newStatement().Dot(name) }
func (g *Group) Dot(name string) *Statement { return g.item(Dot(name)) }
func (s *Statement) Dot(name string) *Statement {
	*s = append(*s,
		token{typ: delimiterToken, content: "."},
		token{typ: identifierToken, content: name},
	)
	return s
}

// ID renders an identifier.
func ID(name string) *Statement                { return newToken(identifierToken, name) }
func (g *Group) ID(name string) *Statement     { return g.item(ID(name)) }
func (s *Statement) ID(name string) *Statement { return s.addToken(identifierToken, name) }

// Qual renders a qualified identifier. Imports are automatically added when
// used with a File. If the path matches the local path, the package name is
// omitted. If package names conflict they are automatically renamed. Note that
// it is not possible to reliably determine the package name given an arbitrary
// package path, so a sensible name is guessed from the path and added as an
// alias. The names of all standard library packages are known so these do not
// need to be aliased. If more control is needed of the aliases, see
// [File.ImportName] or [File.ImportAlias].
func Qual(path, name string) *Statement            { return newStatement().Qual(path, name) }
func (g *Group) Qual(path, name string) *Statement { return g.item(Qual(path, name)) }
func (s *Statement) Qual(path, name string) *Statement {
	*s = append(*s, qualToken{path: path, name: name})
	return s
}

// Line inserts a blank line.
func Line() *Statement                { return newToken(layoutToken, "\n") }
func (g *Group) Line() *Statement     { return g.item(Line()) }
func (s *Statement) Line() *Statement { return s.addToken(layoutToken, "\n") }
