package jen

// This file defines all keyword and identifier token methods. Each is a thin
// wrapper around the token helpers defined in token.go.

// --- Predeclared identifiers (types) ---

// Bool renders the bool identifier.
func Bool() *Statement                { return newToken(identifierToken, "bool") }
func (g *Group) Bool() *Statement     { return g.item(Bool()) }
func (s *Statement) Bool() *Statement { return s.addToken(identifierToken, "bool") }

// Byte renders the byte identifier.
func Byte() *Statement                { return newToken(identifierToken, "byte") }
func (g *Group) Byte() *Statement     { return g.item(Byte()) }
func (s *Statement) Byte() *Statement { return s.addToken(identifierToken, "byte") }

// Complex64 renders the complex64 identifier.
func Complex64() *Statement                { return newToken(identifierToken, "complex64") }
func (g *Group) Complex64() *Statement     { return g.item(Complex64()) }
func (s *Statement) Complex64() *Statement { return s.addToken(identifierToken, "complex64") }

// Complex128 renders the complex128 identifier.
func Complex128() *Statement                { return newToken(identifierToken, "complex128") }
func (g *Group) Complex128() *Statement     { return g.item(Complex128()) }
func (s *Statement) Complex128() *Statement { return s.addToken(identifierToken, "complex128") }

// Error renders the error identifier.
func Error() *Statement                { return newToken(identifierToken, "error") }
func (g *Group) Error() *Statement     { return g.item(Error()) }
func (s *Statement) Error() *Statement { return s.addToken(identifierToken, "error") }

// Float32 renders the float32 identifier.
func Float32() *Statement                { return newToken(identifierToken, "float32") }
func (g *Group) Float32() *Statement     { return g.item(Float32()) }
func (s *Statement) Float32() *Statement { return s.addToken(identifierToken, "float32") }

// Float64 renders the float64 identifier.
func Float64() *Statement                { return newToken(identifierToken, "float64") }
func (g *Group) Float64() *Statement     { return g.item(Float64()) }
func (s *Statement) Float64() *Statement { return s.addToken(identifierToken, "float64") }

// Int renders the int identifier.
func Int() *Statement                { return newToken(identifierToken, "int") }
func (g *Group) Int() *Statement     { return g.item(Int()) }
func (s *Statement) Int() *Statement { return s.addToken(identifierToken, "int") }

// Int8 renders the int8 identifier.
func Int8() *Statement                { return newToken(identifierToken, "int8") }
func (g *Group) Int8() *Statement     { return g.item(Int8()) }
func (s *Statement) Int8() *Statement { return s.addToken(identifierToken, "int8") }

// Int16 renders the int16 identifier.
func Int16() *Statement                { return newToken(identifierToken, "int16") }
func (g *Group) Int16() *Statement     { return g.item(Int16()) }
func (s *Statement) Int16() *Statement { return s.addToken(identifierToken, "int16") }

// Int32 renders the int32 identifier.
func Int32() *Statement                { return newToken(identifierToken, "int32") }
func (g *Group) Int32() *Statement     { return g.item(Int32()) }
func (s *Statement) Int32() *Statement { return s.addToken(identifierToken, "int32") }

// Int64 renders the int64 identifier.
func Int64() *Statement                { return newToken(identifierToken, "int64") }
func (g *Group) Int64() *Statement     { return g.item(Int64()) }
func (s *Statement) Int64() *Statement { return s.addToken(identifierToken, "int64") }

// Rune renders the rune identifier.
func Rune() *Statement                { return newToken(identifierToken, "rune") }
func (g *Group) Rune() *Statement     { return g.item(Rune()) }
func (s *Statement) Rune() *Statement { return s.addToken(identifierToken, "rune") }

// String renders the string identifier.
func String() *Statement                { return newToken(identifierToken, "string") }
func (g *Group) String() *Statement     { return g.item(String()) }
func (s *Statement) String() *Statement { return s.addToken(identifierToken, "string") }

// Uint renders the uint identifier.
func Uint() *Statement                { return newToken(identifierToken, "uint") }
func (g *Group) Uint() *Statement     { return g.item(Uint()) }
func (s *Statement) Uint() *Statement { return s.addToken(identifierToken, "uint") }

// Uint8 renders the uint8 identifier.
func Uint8() *Statement                { return newToken(identifierToken, "uint8") }
func (g *Group) Uint8() *Statement     { return g.item(Uint8()) }
func (s *Statement) Uint8() *Statement { return s.addToken(identifierToken, "uint8") }

// Uint16 renders the uint16 identifier.
func Uint16() *Statement                { return newToken(identifierToken, "uint16") }
func (g *Group) Uint16() *Statement     { return g.item(Uint16()) }
func (s *Statement) Uint16() *Statement { return s.addToken(identifierToken, "uint16") }

// Uint32 renders the uint32 identifier.
func Uint32() *Statement                { return newToken(identifierToken, "uint32") }
func (g *Group) Uint32() *Statement     { return g.item(Uint32()) }
func (s *Statement) Uint32() *Statement { return s.addToken(identifierToken, "uint32") }

// Uint64 renders the uint64 identifier.
func Uint64() *Statement                { return newToken(identifierToken, "uint64") }
func (g *Group) Uint64() *Statement     { return g.item(Uint64()) }
func (s *Statement) Uint64() *Statement { return s.addToken(identifierToken, "uint64") }

// Uintptr renders the uintptr identifier.
func Uintptr() *Statement                { return newToken(identifierToken, "uintptr") }
func (g *Group) Uintptr() *Statement     { return g.item(Uintptr()) }
func (s *Statement) Uintptr() *Statement { return s.addToken(identifierToken, "uintptr") }

// --- Predeclared constants and zero value ---

// True renders the true identifier.
func True() *Statement                { return newToken(identifierToken, "true") }
func (g *Group) True() *Statement     { return g.item(True()) }
func (s *Statement) True() *Statement { return s.addToken(identifierToken, "true") }

// False renders the false identifier.
func False() *Statement                { return newToken(identifierToken, "false") }
func (g *Group) False() *Statement     { return g.item(False()) }
func (s *Statement) False() *Statement { return s.addToken(identifierToken, "false") }

// Iota renders the iota identifier.
func Iota() *Statement                { return newToken(identifierToken, "iota") }
func (g *Group) Iota() *Statement     { return g.item(Iota()) }
func (s *Statement) Iota() *Statement { return s.addToken(identifierToken, "iota") }

// Nil renders the nil identifier.
func Nil() *Statement                { return newToken(identifierToken, "nil") }
func (g *Group) Nil() *Statement     { return g.item(Nil()) }
func (s *Statement) Nil() *Statement { return s.addToken(identifierToken, "nil") }

// --- Common variable names ---

// Err renders the err identifier.
func Err() *Statement                { return newToken(identifierToken, "err") }
func (g *Group) Err() *Statement     { return g.item(Err()) }
func (s *Statement) Err() *Statement { return s.addToken(identifierToken, "err") }

// --- Generics identifiers ---

// Any renders the any identifier.
func Any() *Statement                { return newToken(identifierToken, "any") }
func (g *Group) Any() *Statement     { return g.item(Any()) }
func (s *Statement) Any() *Statement { return s.addToken(identifierToken, "any") }

// Comparable renders the comparable identifier.
func Comparable() *Statement                { return newToken(identifierToken, "comparable") }
func (g *Group) Comparable() *Statement     { return g.item(Comparable()) }
func (s *Statement) Comparable() *Statement { return s.addToken(identifierToken, "comparable") }

// --- Keywords ---

// Break renders the break keyword.
func Break() *Statement                { return newToken(keywordToken, "break") }
func (g *Group) Break() *Statement     { return g.item(Break()) }
func (s *Statement) Break() *Statement { return s.addToken(keywordToken, "break") }

// Chan renders the chan keyword.
func Chan() *Statement                { return newToken(keywordToken, "chan") }
func (g *Group) Chan() *Statement     { return g.item(Chan()) }
func (s *Statement) Chan() *Statement { return s.addToken(keywordToken, "chan") }

// Const renders the const keyword.
func Const() *Statement                { return newToken(keywordToken, "const") }
func (g *Group) Const() *Statement     { return g.item(Const()) }
func (s *Statement) Const() *Statement { return s.addToken(keywordToken, "const") }

// Continue renders the continue keyword.
func Continue() *Statement                { return newToken(keywordToken, "continue") }
func (g *Group) Continue() *Statement     { return g.item(Continue()) }
func (s *Statement) Continue() *Statement { return s.addToken(keywordToken, "continue") }

// Defer renders the defer keyword.
func Defer() *Statement                { return newToken(keywordToken, "defer") }
func (g *Group) Defer() *Statement     { return g.item(Defer()) }
func (s *Statement) Defer() *Statement { return s.addToken(keywordToken, "defer") }

// Else renders the else keyword.
func Else() *Statement                { return newToken(keywordToken, "else") }
func (g *Group) Else() *Statement     { return g.item(Else()) }
func (s *Statement) Else() *Statement { return s.addToken(keywordToken, "else") }

// Fallthrough renders the fallthrough keyword.
func Fallthrough() *Statement                { return newToken(keywordToken, "fallthrough") }
func (g *Group) Fallthrough() *Statement     { return g.item(Fallthrough()) }
func (s *Statement) Fallthrough() *Statement { return s.addToken(keywordToken, "fallthrough") }

// Func renders the func keyword.
func Func() *Statement                { return newToken(keywordToken, "func") }
func (g *Group) Func() *Statement     { return g.item(Func()) }
func (s *Statement) Func() *Statement { return s.addToken(keywordToken, "func") }

// Go renders the go keyword.
func Go() *Statement                { return newToken(keywordToken, "go") }
func (g *Group) Go() *Statement     { return g.item(Go()) }
func (s *Statement) Go() *Statement { return s.addToken(keywordToken, "go") }

// Goto renders the goto keyword.
func Goto() *Statement                { return newToken(keywordToken, "goto") }
func (g *Group) Goto() *Statement     { return g.item(Goto()) }
func (s *Statement) Goto() *Statement { return s.addToken(keywordToken, "goto") }

// Range renders the range keyword.
func Range() *Statement                { return newToken(keywordToken, "range") }
func (g *Group) Range() *Statement     { return g.item(Range()) }
func (s *Statement) Range() *Statement { return s.addToken(keywordToken, "range") }

// Select renders the select keyword.
func Select() *Statement                { return newToken(keywordToken, "select") }
func (g *Group) Select() *Statement     { return g.item(Select()) }
func (s *Statement) Select() *Statement { return s.addToken(keywordToken, "select") }

// Type renders the type keyword.
func Type() *Statement                { return newToken(keywordToken, "type") }
func (g *Group) Type() *Statement     { return g.item(Type()) }
func (s *Statement) Type() *Statement { return s.addToken(keywordToken, "type") }

// Var renders the var keyword.
func Var() *Statement                { return newToken(keywordToken, "var") }
func (g *Group) Var() *Statement     { return g.item(Var()) }
func (s *Statement) Var() *Statement { return s.addToken(keywordToken, "var") }
