// Package jen is a code generator for Go.
package jen

import (
	"io"
)

// Code represents an item of code that can be rendered.
type Code interface {
	render(f *File, w io.Writer, s *Statement) error
	isNull(f *File) bool
}

var newLine = []byte("\n")
