package jen

// reserved is the canonical set of Go reserved words: keywords, predeclared
// identifiers, and common names that should not be used as import aliases.
// This is the single source of truth — IsReservedWord, import alias validation,
// and the generated keyword/identifier methods all derive from this list.
var reserved = map[string]bool{
	// keywords (https://go.dev/ref/spec#Keywords)
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true, "for": true,
	"func": true, "go": true, "goto": true, "if": true, "import": true,
	"interface": true, "map": true, "package": true, "range": true, "return": true,
	"select": true, "struct": true, "switch": true, "type": true, "var": true,
	// predeclared types
	"any": true, "bool": true, "byte": true, "comparable": true,
	"complex64": true, "complex128": true, "error": true,
	"float32": true, "float64": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true,
	// predeclared constants
	"true": true, "false": true, "iota": true,
	// predeclared zero value
	"nil": true,
	// predeclared functions
	"append": true, "cap": true, "clear": true, "close": true, "complex": true,
	"copy": true, "delete": true, "imag": true, "len": true, "make": true,
	"max": true, "min": true, "new": true, "panic": true, "print": true,
	"println": true, "real": true, "recover": true,
	// common variables
	"err": true,
}

// IsReservedWord returns true if alias is a Go keyword, predeclared identifier,
// or common variable name that should not be used as an import alias.
func IsReservedWord(alias string) bool {
	return reserved[alias]
}
