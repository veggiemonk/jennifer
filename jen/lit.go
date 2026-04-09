package jen

// Lit renders a literal. Lit supports only built-in types (bool, string, int, complex128, float64,
// float32, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Passing any other type will panic.
func Lit(v any) *Statement            { return newStatement().Lit(v) }
func (g *Group) Lit(v any) *Statement { return g.item(Lit(v)) }
func (s *Statement) Lit(v any) *Statement {
	*s = append(*s, token{typ: literalToken, content: v})
	return s
}

// LitFunc renders a literal. LitFunc generates the value to render by executing the provided
// function. LitFunc supports only built-in types (bool, string, int, complex128, float64, float32,
// int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr and complex64).
// Returning any other type will panic.
func LitFunc(f func() any) *Statement            { return newStatement().LitFunc(f) }
func (g *Group) LitFunc(f func() any) *Statement { return g.item(LitFunc(f)) }
func (s *Statement) LitFunc(f func() any) *Statement {
	*s = append(*s, token{typ: literalToken, content: f()})
	return s
}

// LitRune renders a rune literal.
func LitRune(v rune) *Statement            { return newStatement().LitRune(v) }
func (g *Group) LitRune(v rune) *Statement { return g.item(LitRune(v)) }
func (s *Statement) LitRune(v rune) *Statement {
	*s = append(*s, token{typ: literalRuneToken, content: v})
	return s
}

// LitRuneFunc renders a rune literal. LitRuneFunc generates the value to
// render by executing the provided function.
func LitRuneFunc(f func() rune) *Statement            { return newStatement().LitRuneFunc(f) }
func (g *Group) LitRuneFunc(f func() rune) *Statement { return g.item(LitRuneFunc(f)) }
func (s *Statement) LitRuneFunc(f func() rune) *Statement {
	*s = append(*s, token{typ: literalRuneToken, content: f()})
	return s
}

// LitByte renders a byte literal.
func LitByte(v byte) *Statement            { return newStatement().LitByte(v) }
func (g *Group) LitByte(v byte) *Statement { return g.item(LitByte(v)) }
func (s *Statement) LitByte(v byte) *Statement {
	*s = append(*s, token{typ: literalByteToken, content: v})
	return s
}

// LitByteFunc renders a byte literal. LitByteFunc generates the value to
// render by executing the provided function.
func LitByteFunc(f func() byte) *Statement            { return newStatement().LitByteFunc(f) }
func (g *Group) LitByteFunc(f func() byte) *Statement { return g.item(LitByteFunc(f)) }
func (s *Statement) LitByteFunc(f func() byte) *Statement {
	*s = append(*s, token{typ: literalByteToken, content: f()})
	return s
}
