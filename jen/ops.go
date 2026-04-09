package jen

// This file defines all group-based operations. Each operation is a thin
// wrapper around the group helpers defined in group.go. No code generation
// needed — adding a new Go language construct is just a few lines here.

// --- Group configurations ---

var (
	parensOp    = groupOp{"(", ")", "", false}
	listOp      = groupOp{"", "", ",", false}
	valuesOp    = groupOp{"{", "}", ",", false}
	indexOp     = groupOp{"[", "]", ":", false}
	blockOp     = groupOp{"{", "}", "", true}
	defsOp      = groupOp{"(", ")", "", true}
	callOp      = groupOp{"(", ")", ",", false}
	paramsOp    = groupOp{"(", ")", ",", false}
	assertOp    = groupOp{".(", ")", "", false}
	mapOp       = groupOp{"map[", "]", "", false}
	ifOp        = groupOp{"if ", "", ";", false}
	returnOp    = groupOp{"return ", "", ",", false}
	forOp       = groupOp{"for ", "", ";", false}
	switchOp    = groupOp{"switch ", "", ";", false}
	interfaceOp = groupOp{"interface{", "}", "", true}
	structOp    = groupOp{"struct{", "}", "", true}
	caseOp      = groupOp{"case ", ":", ",", false}
	defaultOp   = groupOp{"default", ":", "", false}
	typesOp     = groupOp{"[", "]", ",", false}
	unionOp     = groupOp{"", "", "|", false}
	// built-in functions
	appendOp  = groupOp{"append(", ")", ",", false}
	capOp     = groupOp{"cap(", ")", ",", false}
	closeOpG  = groupOp{"close(", ")", ",", false}
	clearOp   = groupOp{"clear(", ")", ",", false}
	minOp     = groupOp{"min(", ")", ",", false}
	maxOp     = groupOp{"max(", ")", ",", false}
	complexOp = groupOp{"complex(", ")", ",", false}
	copyOp    = groupOp{"copy(", ")", ",", false}
	deleteOp  = groupOp{"delete(", ")", ",", false}
	imagOp    = groupOp{"imag(", ")", ",", false}
	lenOp     = groupOp{"len(", ")", ",", false}
	makeOp    = groupOp{"make(", ")", ",", false}
	newOpG    = groupOp{"new(", ")", ",", false}
	panicOp   = groupOp{"panic(", ")", ",", false}
	printOp   = groupOp{"print(", ")", ",", false}
	printlnOp = groupOp{"println(", ")", ",", false}
	realOp    = groupOp{"real(", ")", ",", false}
	recoverOp = groupOp{"recover(", ")", ",", false}
)

// --- Language constructs ---

// Parens renders a single item in parenthesis. Use for type conversion or to specify evaluation order.
func Parens(item Code) *Statement                { return newGroup("parens", parensOp, []Code{item}) }
func (g *Group) Parens(item Code) *Statement     { return g.item(Parens(item)) }
func (s *Statement) Parens(item Code) *Statement { return s.addGroup("parens", parensOp, []Code{item}) }

// List renders a comma separated list. Use for multiple return functions.
func List(items ...Code) *Statement                { return newGroup("list", listOp, items) }
func (g *Group) List(items ...Code) *Statement     { return g.item(List(items...)) }
func (s *Statement) List(items ...Code) *Statement { return s.addGroup("list", listOp, items) }

// ListFunc renders a comma separated list. Use for multiple return functions.
func ListFunc(f func(*Group)) *Statement                { return newGroupFunc("list", listOp, f) }
func (g *Group) ListFunc(f func(*Group)) *Statement     { return g.item(ListFunc(f)) }
func (s *Statement) ListFunc(f func(*Group)) *Statement { return s.addGroupFunc("list", listOp, f) }

// Values renders a comma separated list enclosed by curly braces. Use for slice or composite literals.
func Values(values ...Code) *Statement                { return newGroup("values", valuesOp, values) }
func (g *Group) Values(values ...Code) *Statement     { return g.item(Values(values...)) }
func (s *Statement) Values(values ...Code) *Statement { return s.addGroup("values", valuesOp, values) }

// ValuesFunc renders a comma separated list enclosed by curly braces. Use for slice or composite literals.
func ValuesFunc(f func(*Group)) *Statement            { return newGroupFunc("values", valuesOp, f) }
func (g *Group) ValuesFunc(f func(*Group)) *Statement { return g.item(ValuesFunc(f)) }
func (s *Statement) ValuesFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("values", valuesOp, f)
}

// Index renders a colon separated list enclosed by square brackets. Use for array / slice indexes and definitions.
func Index(items ...Code) *Statement                { return newGroup("index", indexOp, items) }
func (g *Group) Index(items ...Code) *Statement     { return g.item(Index(items...)) }
func (s *Statement) Index(items ...Code) *Statement { return s.addGroup("index", indexOp, items) }

// IndexFunc renders a colon separated list enclosed by square brackets. Use for array / slice indexes and definitions.
func IndexFunc(f func(*Group)) *Statement                { return newGroupFunc("index", indexOp, f) }
func (g *Group) IndexFunc(f func(*Group)) *Statement     { return g.item(IndexFunc(f)) }
func (s *Statement) IndexFunc(f func(*Group)) *Statement { return s.addGroupFunc("index", indexOp, f) }

// Block renders a statement list enclosed by curly braces. Use for code blocks.
// A special case applies when used directly after Case or Default, where the braces
// are omitted. This allows use in switch and select statements.
func Block(statements ...Code) *Statement            { return newGroup("block", blockOp, statements) }
func (g *Group) Block(statements ...Code) *Statement { return g.item(Block(statements...)) }
func (s *Statement) Block(statements ...Code) *Statement {
	return s.addGroup("block", blockOp, statements)
}

// BlockFunc renders a statement list enclosed by curly braces. Use for code blocks.
func BlockFunc(f func(*Group)) *Statement                { return newGroupFunc("block", blockOp, f) }
func (g *Group) BlockFunc(f func(*Group)) *Statement     { return g.item(BlockFunc(f)) }
func (s *Statement) BlockFunc(f func(*Group)) *Statement { return s.addGroupFunc("block", blockOp, f) }

// Defs renders a statement list enclosed in parenthesis. Use for definition lists.
func Defs(definitions ...Code) *Statement            { return newGroup("defs", defsOp, definitions) }
func (g *Group) Defs(definitions ...Code) *Statement { return g.item(Defs(definitions...)) }
func (s *Statement) Defs(definitions ...Code) *Statement {
	return s.addGroup("defs", defsOp, definitions)
}

// DefsFunc renders a statement list enclosed in parenthesis. Use for definition lists.
func DefsFunc(f func(*Group)) *Statement                { return newGroupFunc("defs", defsOp, f) }
func (g *Group) DefsFunc(f func(*Group)) *Statement     { return g.item(DefsFunc(f)) }
func (s *Statement) DefsFunc(f func(*Group)) *Statement { return s.addGroupFunc("defs", defsOp, f) }

// Call renders a comma separated list enclosed by parenthesis. Use for function calls.
func Call(params ...Code) *Statement                { return newGroup("call", callOp, params) }
func (g *Group) Call(params ...Code) *Statement     { return g.item(Call(params...)) }
func (s *Statement) Call(params ...Code) *Statement { return s.addGroup("call", callOp, params) }

// CallFunc renders a comma separated list enclosed by parenthesis. Use for function calls.
func CallFunc(f func(*Group)) *Statement                { return newGroupFunc("call", callOp, f) }
func (g *Group) CallFunc(f func(*Group)) *Statement     { return g.item(CallFunc(f)) }
func (s *Statement) CallFunc(f func(*Group)) *Statement { return s.addGroupFunc("call", callOp, f) }

// Params renders a comma separated list enclosed by parenthesis. Use for function parameters and method receivers.
func Params(params ...Code) *Statement                { return newGroup("params", paramsOp, params) }
func (g *Group) Params(params ...Code) *Statement     { return g.item(Params(params...)) }
func (s *Statement) Params(params ...Code) *Statement { return s.addGroup("params", paramsOp, params) }

// ParamsFunc renders a comma separated list enclosed by parenthesis. Use for function parameters and method receivers.
func ParamsFunc(f func(*Group)) *Statement            { return newGroupFunc("params", paramsOp, f) }
func (g *Group) ParamsFunc(f func(*Group)) *Statement { return g.item(ParamsFunc(f)) }
func (s *Statement) ParamsFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("params", paramsOp, f)
}

// Assert renders a period followed by a single item enclosed by parenthesis. Use for type assertions.
func Assert(typ Code) *Statement                { return newGroup("assert", assertOp, []Code{typ}) }
func (g *Group) Assert(typ Code) *Statement     { return g.item(Assert(typ)) }
func (s *Statement) Assert(typ Code) *Statement { return s.addGroup("assert", assertOp, []Code{typ}) }

// Map renders the keyword followed by a single item enclosed by square brackets. Use for map definitions.
func Map(typ Code) *Statement                { return newGroup("map", mapOp, []Code{typ}) }
func (g *Group) Map(typ Code) *Statement     { return g.item(Map(typ)) }
func (s *Statement) Map(typ Code) *Statement { return s.addGroup("map", mapOp, []Code{typ}) }

// If renders the keyword followed by a semicolon separated list.
func If(conditions ...Code) *Statement                { return newGroup("if", ifOp, conditions) }
func (g *Group) If(conditions ...Code) *Statement     { return g.item(If(conditions...)) }
func (s *Statement) If(conditions ...Code) *Statement { return s.addGroup("if", ifOp, conditions) }

// IfFunc renders the keyword followed by a semicolon separated list.
func IfFunc(f func(*Group)) *Statement                { return newGroupFunc("if", ifOp, f) }
func (g *Group) IfFunc(f func(*Group)) *Statement     { return g.item(IfFunc(f)) }
func (s *Statement) IfFunc(f func(*Group)) *Statement { return s.addGroupFunc("if", ifOp, f) }

// Return renders the keyword followed by a comma separated list.
func Return(results ...Code) *Statement            { return newGroup("return", returnOp, results) }
func (g *Group) Return(results ...Code) *Statement { return g.item(Return(results...)) }
func (s *Statement) Return(results ...Code) *Statement {
	return s.addGroup("return", returnOp, results)
}

// ReturnFunc renders the keyword followed by a comma separated list.
func ReturnFunc(f func(*Group)) *Statement            { return newGroupFunc("return", returnOp, f) }
func (g *Group) ReturnFunc(f func(*Group)) *Statement { return g.item(ReturnFunc(f)) }
func (s *Statement) ReturnFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("return", returnOp, f)
}

// For renders the keyword followed by a semicolon separated list.
func For(conditions ...Code) *Statement                { return newGroup("for", forOp, conditions) }
func (g *Group) For(conditions ...Code) *Statement     { return g.item(For(conditions...)) }
func (s *Statement) For(conditions ...Code) *Statement { return s.addGroup("for", forOp, conditions) }

// ForFunc renders the keyword followed by a semicolon separated list.
func ForFunc(f func(*Group)) *Statement                { return newGroupFunc("for", forOp, f) }
func (g *Group) ForFunc(f func(*Group)) *Statement     { return g.item(ForFunc(f)) }
func (s *Statement) ForFunc(f func(*Group)) *Statement { return s.addGroupFunc("for", forOp, f) }

// Switch renders the keyword followed by a semicolon separated list.
func Switch(conditions ...Code) *Statement            { return newGroup("switch", switchOp, conditions) }
func (g *Group) Switch(conditions ...Code) *Statement { return g.item(Switch(conditions...)) }
func (s *Statement) Switch(conditions ...Code) *Statement {
	return s.addGroup("switch", switchOp, conditions)
}

// SwitchFunc renders the keyword followed by a semicolon separated list.
func SwitchFunc(f func(*Group)) *Statement            { return newGroupFunc("switch", switchOp, f) }
func (g *Group) SwitchFunc(f func(*Group)) *Statement { return g.item(SwitchFunc(f)) }
func (s *Statement) SwitchFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("switch", switchOp, f)
}

// Interface renders the keyword followed by a method list enclosed by curly braces.
func Interface(methods ...Code) *Statement            { return newGroup("interface", interfaceOp, methods) }
func (g *Group) Interface(methods ...Code) *Statement { return g.item(Interface(methods...)) }
func (s *Statement) Interface(methods ...Code) *Statement {
	return s.addGroup("interface", interfaceOp, methods)
}

// InterfaceFunc renders the keyword followed by a method list enclosed by curly braces.
func InterfaceFunc(f func(*Group)) *Statement            { return newGroupFunc("interface", interfaceOp, f) }
func (g *Group) InterfaceFunc(f func(*Group)) *Statement { return g.item(InterfaceFunc(f)) }
func (s *Statement) InterfaceFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("interface", interfaceOp, f)
}

// Struct renders the keyword followed by a field list enclosed by curly braces.
func Struct(fields ...Code) *Statement                { return newGroup("struct", structOp, fields) }
func (g *Group) Struct(fields ...Code) *Statement     { return g.item(Struct(fields...)) }
func (s *Statement) Struct(fields ...Code) *Statement { return s.addGroup("struct", structOp, fields) }

// StructFunc renders the keyword followed by a field list enclosed by curly braces.
func StructFunc(f func(*Group)) *Statement            { return newGroupFunc("struct", structOp, f) }
func (g *Group) StructFunc(f func(*Group)) *Statement { return g.item(StructFunc(f)) }
func (s *Statement) StructFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("struct", structOp, f)
}

// Case renders the keyword followed by a comma separated list.
func Case(cases ...Code) *Statement                { return newGroup("case", caseOp, cases) }
func (g *Group) Case(cases ...Code) *Statement     { return g.item(Case(cases...)) }
func (s *Statement) Case(cases ...Code) *Statement { return s.addGroup("case", caseOp, cases) }

// CaseFunc renders the keyword followed by a comma separated list.
func CaseFunc(f func(*Group)) *Statement                { return newGroupFunc("case", caseOp, f) }
func (g *Group) CaseFunc(f func(*Group)) *Statement     { return g.item(CaseFunc(f)) }
func (s *Statement) CaseFunc(f func(*Group)) *Statement { return s.addGroupFunc("case", caseOp, f) }

// Default renders the default keyword followed by a colon.
func Default() *Statement                { return newGroup("default", defaultOp, nil) }
func (g *Group) Default() *Statement     { return g.item(Default()) }
func (s *Statement) Default() *Statement { return s.addGroup("default", defaultOp, nil) }

// Types renders a comma separated list enclosed by square brackets. Use for type parameters and constraints.
func Types(types ...Code) *Statement                { return newGroup("types", typesOp, types) }
func (g *Group) Types(types ...Code) *Statement     { return g.item(Types(types...)) }
func (s *Statement) Types(types ...Code) *Statement { return s.addGroup("types", typesOp, types) }

// TypesFunc renders a comma separated list enclosed by square brackets. Use for type parameters and constraints.
func TypesFunc(f func(*Group)) *Statement                { return newGroupFunc("types", typesOp, f) }
func (g *Group) TypesFunc(f func(*Group)) *Statement     { return g.item(TypesFunc(f)) }
func (s *Statement) TypesFunc(f func(*Group)) *Statement { return s.addGroupFunc("types", typesOp, f) }

// Union renders a pipe separated list. Use for union type constraints.
func Union(types ...Code) *Statement                { return newGroup("union", unionOp, types) }
func (g *Group) Union(types ...Code) *Statement     { return g.item(Union(types...)) }
func (s *Statement) Union(types ...Code) *Statement { return s.addGroup("union", unionOp, types) }

// UnionFunc renders a pipe separated list. Use for union type constraints.
func UnionFunc(f func(*Group)) *Statement                { return newGroupFunc("union", unionOp, f) }
func (g *Group) UnionFunc(f func(*Group)) *Statement     { return g.item(UnionFunc(f)) }
func (s *Statement) UnionFunc(f func(*Group)) *Statement { return s.addGroupFunc("union", unionOp, f) }

// --- Built-in functions ---

// Append renders the append built-in function.
func Append(args ...Code) *Statement                { return newGroup("append", appendOp, args) }
func (g *Group) Append(args ...Code) *Statement     { return g.item(Append(args...)) }
func (s *Statement) Append(args ...Code) *Statement { return s.addGroup("append", appendOp, args) }

// AppendFunc renders the append built-in function.
func AppendFunc(f func(*Group)) *Statement            { return newGroupFunc("append", appendOp, f) }
func (g *Group) AppendFunc(f func(*Group)) *Statement { return g.item(AppendFunc(f)) }
func (s *Statement) AppendFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("append", appendOp, f)
}

// Cap renders the cap built-in function.
func Cap(v Code) *Statement                { return newGroup("cap", capOp, []Code{v}) }
func (g *Group) Cap(v Code) *Statement     { return g.item(Cap(v)) }
func (s *Statement) Cap(v Code) *Statement { return s.addGroup("cap", capOp, []Code{v}) }

// Close renders the close built-in function.
func Close(c Code) *Statement                { return newGroup("close", closeOpG, []Code{c}) }
func (g *Group) Close(c Code) *Statement     { return g.item(Close(c)) }
func (s *Statement) Close(c Code) *Statement { return s.addGroup("close", closeOpG, []Code{c}) }

// Clear renders the clear built-in function.
func Clear(c Code) *Statement                { return newGroup("clear", clearOp, []Code{c}) }
func (g *Group) Clear(c Code) *Statement     { return g.item(Clear(c)) }
func (s *Statement) Clear(c Code) *Statement { return s.addGroup("clear", clearOp, []Code{c}) }

// Min renders the min built-in function.
func Min(args ...Code) *Statement                { return newGroup("min", minOp, args) }
func (g *Group) Min(args ...Code) *Statement     { return g.item(Min(args...)) }
func (s *Statement) Min(args ...Code) *Statement { return s.addGroup("min", minOp, args) }

// MinFunc renders the min built-in function.
func MinFunc(f func(*Group)) *Statement                { return newGroupFunc("min", minOp, f) }
func (g *Group) MinFunc(f func(*Group)) *Statement     { return g.item(MinFunc(f)) }
func (s *Statement) MinFunc(f func(*Group)) *Statement { return s.addGroupFunc("min", minOp, f) }

// Max renders the max built-in function.
func Max(args ...Code) *Statement                { return newGroup("max", maxOp, args) }
func (g *Group) Max(args ...Code) *Statement     { return g.item(Max(args...)) }
func (s *Statement) Max(args ...Code) *Statement { return s.addGroup("max", maxOp, args) }

// MaxFunc renders the max built-in function.
func MaxFunc(f func(*Group)) *Statement                { return newGroupFunc("max", maxOp, f) }
func (g *Group) MaxFunc(f func(*Group)) *Statement     { return g.item(MaxFunc(f)) }
func (s *Statement) MaxFunc(f func(*Group)) *Statement { return s.addGroupFunc("max", maxOp, f) }

// Complex renders the complex built-in function.
func Complex(r Code, i Code) *Statement            { return newGroup("complex", complexOp, []Code{r, i}) }
func (g *Group) Complex(r Code, i Code) *Statement { return g.item(Complex(r, i)) }
func (s *Statement) Complex(r Code, i Code) *Statement {
	return s.addGroup("complex", complexOp, []Code{r, i})
}

// Copy renders the copy built-in function.
func Copy(dst Code, src Code) *Statement            { return newGroup("copy", copyOp, []Code{dst, src}) }
func (g *Group) Copy(dst Code, src Code) *Statement { return g.item(Copy(dst, src)) }
func (s *Statement) Copy(dst Code, src Code) *Statement {
	return s.addGroup("copy", copyOp, []Code{dst, src})
}

// Delete renders the delete built-in function.
func Delete(m Code, key Code) *Statement            { return newGroup("delete", deleteOp, []Code{m, key}) }
func (g *Group) Delete(m Code, key Code) *Statement { return g.item(Delete(m, key)) }
func (s *Statement) Delete(m Code, key Code) *Statement {
	return s.addGroup("delete", deleteOp, []Code{m, key})
}

// Imag renders the imag built-in function.
func Imag(c Code) *Statement                { return newGroup("imag", imagOp, []Code{c}) }
func (g *Group) Imag(c Code) *Statement     { return g.item(Imag(c)) }
func (s *Statement) Imag(c Code) *Statement { return s.addGroup("imag", imagOp, []Code{c}) }

// Len renders the len built-in function.
func Len(v Code) *Statement                { return newGroup("len", lenOp, []Code{v}) }
func (g *Group) Len(v Code) *Statement     { return g.item(Len(v)) }
func (s *Statement) Len(v Code) *Statement { return s.addGroup("len", lenOp, []Code{v}) }

// Make renders the make built-in function. The final parameter of the make function
// is optional, so it is represented by a variadic parameter list.
func Make(args ...Code) *Statement                { return newGroup("make", makeOp, args) }
func (g *Group) Make(args ...Code) *Statement     { return g.item(Make(args...)) }
func (s *Statement) Make(args ...Code) *Statement { return s.addGroup("make", makeOp, args) }

// New renders the new built-in function.
func New(typ Code) *Statement                { return newGroup("new", newOpG, []Code{typ}) }
func (g *Group) New(typ Code) *Statement     { return g.item(New(typ)) }
func (s *Statement) New(typ Code) *Statement { return s.addGroup("new", newOpG, []Code{typ}) }

// Panic renders the panic built-in function.
func Panic(v Code) *Statement                { return newGroup("panic", panicOp, []Code{v}) }
func (g *Group) Panic(v Code) *Statement     { return g.item(Panic(v)) }
func (s *Statement) Panic(v Code) *Statement { return s.addGroup("panic", panicOp, []Code{v}) }

// Print renders the print built-in function.
func Print(args ...Code) *Statement                { return newGroup("print", printOp, args) }
func (g *Group) Print(args ...Code) *Statement     { return g.item(Print(args...)) }
func (s *Statement) Print(args ...Code) *Statement { return s.addGroup("print", printOp, args) }

// PrintFunc renders the print built-in function.
func PrintFunc(f func(*Group)) *Statement                { return newGroupFunc("print", printOp, f) }
func (g *Group) PrintFunc(f func(*Group)) *Statement     { return g.item(PrintFunc(f)) }
func (s *Statement) PrintFunc(f func(*Group)) *Statement { return s.addGroupFunc("print", printOp, f) }

// Println renders the println built-in function.
func Println(args ...Code) *Statement                { return newGroup("println", printlnOp, args) }
func (g *Group) Println(args ...Code) *Statement     { return g.item(Println(args...)) }
func (s *Statement) Println(args ...Code) *Statement { return s.addGroup("println", printlnOp, args) }

// PrintlnFunc renders the println built-in function.
func PrintlnFunc(f func(*Group)) *Statement            { return newGroupFunc("println", printlnOp, f) }
func (g *Group) PrintlnFunc(f func(*Group)) *Statement { return g.item(PrintlnFunc(f)) }
func (s *Statement) PrintlnFunc(f func(*Group)) *Statement {
	return s.addGroupFunc("println", printlnOp, f)
}

// Real renders the real built-in function.
func Real(c Code) *Statement                { return newGroup("real", realOp, []Code{c}) }
func (g *Group) Real(c Code) *Statement     { return g.item(Real(c)) }
func (s *Statement) Real(c Code) *Statement { return s.addGroup("real", realOp, []Code{c}) }

// Recover renders the recover built-in function.
func Recover() *Statement                { return newGroup("recover", recoverOp, nil) }
func (g *Group) Recover() *Statement     { return g.item(Recover()) }
func (s *Statement) Recover() *Statement { return s.addGroup("recover", recoverOp, nil) }
