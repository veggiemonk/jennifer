package jen

import (
	"fmt"
	"io"
	"strings"
)

// Comment adds a comment. If the provided string contains a newline, the
// comment is formatted in multiline style. If the comment string starts
// with "//" or "/*", the automatic formatting is disabled and the string is
// rendered directly.
func Comment(str string) *Statement            { return newStatement().Comment(str) }
func (g *Group) Comment(str string) *Statement { return g.item(Comment(str)) }
func (s *Statement) Comment(str string) *Statement {
	*s = append(*s, comment{comment: str})
	return s
}

// Commentf adds a comment, using a format string and a list of parameters. If
// the provided string contains a newline, the comment is formatted in
// multiline style. If the comment string starts with "//" or "/*", the
// automatic formatting is disabled and the string is rendered directly.
func Commentf(format string, a ...any) *Statement            { return newStatement().Commentf(format, a...) }
func (g *Group) Commentf(format string, a ...any) *Statement { return g.item(Commentf(format, a...)) }
func (s *Statement) Commentf(format string, a ...any) *Statement {
	*s = append(*s, comment{comment: fmt.Sprintf(format, a...)})
	return s
}

type comment struct {
	comment string
}

func (c comment) isNull(f *File) bool { return false }

func (c comment) render(f *File, w io.Writer, s *Statement) error {
	if strings.HasPrefix(c.comment, "//") || strings.HasPrefix(c.comment, "/*") {
		_, err := w.Write([]byte(c.comment))
		return err
	}
	if strings.Contains(c.comment, "\n") {
		if _, err := w.Write([]byte("/*\n")); err != nil {
			return err
		}
	} else {
		if _, err := w.Write([]byte("// ")); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte(c.comment)); err != nil {
		return err
	}
	if strings.Contains(c.comment, "\n") {
		if !strings.HasSuffix(c.comment, "\n") {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte("*/")); err != nil {
			return err
		}
	}
	return nil
}
