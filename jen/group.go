package jen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

// Group represents a list of Code items, separated by tokens with an optional
// open and close token.
type Group struct {
	name      string
	items     []Code
	open      string
	close     string
	separator string
	multi     bool
}

func (g *Group) isNull(f *File) bool {
	if g == nil {
		return true
	}
	if g.open != "" || g.close != "" {
		return false
	}
	return g.isNullItems(f)
}

func (g *Group) isNullItems(f *File) bool {
	for _, c := range g.items {
		if !c.isNull(f) {
			return false
		}
	}
	return true
}

func (g *Group) render(f *File, w io.Writer, s *Statement) error {
	if g.name == "types" && g.isNullItems(f) {
		// Special case for types - if all items are null, don't render the open/close tokens.
		return nil
	}
	open := g.open
	closeTok := g.close
	if g.name == "block" && s != nil {
		// When a block follows a Case or Default group, omit braces (switch/select syntax).
		prev := s.previous(g)
		if grp, ok := prev.(*Group); ok && (grp.name == "case" || grp.name == "default") {
			open = ""
			closeTok = ""
		}
	}
	if open != "" {
		if _, err := w.Write([]byte(open)); err != nil {
			return err
		}
	}
	isNull, err := g.renderItems(f, w)
	if err != nil {
		return err
	}
	if !isNull && g.multi && closeTok != "" {
		// For multi-line blocks with a closing token, insert a newline after the last item
		// (but not if all items were null). This ensures a trailing comment doesn't
		// swallow the closing token.
		sep := "\n"
		if g.separator == "," {
			sep = ",\n"
		}
		if _, err := w.Write([]byte(sep)); err != nil {
			return err
		}
	}
	if closeTok != "" {
		if _, err := w.Write([]byte(closeTok)); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) renderItems(f *File, w io.Writer) (isNull bool, err error) {
	first := true
	for _, code := range g.items {
		if code == nil || code.isNull(f) {
			continue
		}
		if g.name == "values" {
			if _, ok := code.(Dict); ok && len(g.items) > 1 {
				panic("Error in Values: if Dict is used, must be one item only")
			}
		}
		if !first && g.separator != "" {
			if _, err := w.Write([]byte(g.separator)); err != nil {
				return false, err
			}
		}
		if g.multi {
			if _, err := w.Write([]byte("\n")); err != nil {
				return false, err
			}
		}
		if err := code.render(f, w, nil); err != nil {
			return false, err
		}
		first = false
	}
	return first, nil
}

// Render renders the Group to the provided writer.
func (g *Group) Render(writer io.Writer) error {
	return g.RenderWithFile(writer, NewFile(""))
}

// GoString renders the Group for testing. Any error will cause a panic.
func (g *Group) GoString() string {
	buf := bytes.Buffer{}
	if err := g.Render(&buf); err != nil {
		panic(err)
	}
	return buf.String()
}

// RenderWithFile renders the Group to the provided writer, using imports from the provided file.
func (g *Group) RenderWithFile(writer io.Writer, file *File) error {
	buf := &bytes.Buffer{}
	if err := g.render(file, buf, nil); err != nil {
		return fmt.Errorf("rendering group: %w", err)
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting generated source: %w\n%s", err, numberLines(buf.String()))
	}
	if _, err := writer.Write(b); err != nil {
		return err
	}
	return nil
}

// --- Group operation helpers ---

// groupOp defines the rendering configuration for a group-based operation.
type groupOp struct {
	open, close, sep string
	multi            bool
}

// item adds a statement to the group and returns it for chaining.
func (g *Group) item(s *Statement) *Statement {
	g.items = append(g.items, s)
	return s
}

// newGroup creates a new Statement containing a Group with the given configuration.
func newGroup(name string, op groupOp, items []Code) *Statement {
	return newStatement().addGroup(name, op, items)
}

// newGroupFunc creates a new Statement containing a Group populated by a callback.
func newGroupFunc(name string, op groupOp, f func(*Group)) *Statement {
	return newStatement().addGroupFunc(name, op, f)
}

// addGroup appends a Group with the given items to the statement.
func (s *Statement) addGroup(name string, op groupOp, items []Code) *Statement {
	g := &Group{
		name:      name,
		items:     items,
		open:      op.open,
		close:     op.close,
		separator: op.sep,
		multi:     op.multi,
	}
	*s = append(*s, g)
	return s
}

// addGroupFunc appends a Group populated by a callback to the statement.
func (s *Statement) addGroupFunc(name string, op groupOp, f func(*Group)) *Statement {
	g := &Group{
		name:      name,
		open:      op.open,
		close:     op.close,
		separator: op.sep,
		multi:     op.multi,
	}
	f(g)
	*s = append(*s, g)
	return s
}

// --- Options / Custom ---

// Options specifies options for the Custom method.
type Options struct {
	Open      string
	Close     string
	Separator string
	Multi     bool
}

// Custom renders a customized statement list. Pass in options to specify
// multi-line, and tokens for open, close, separator.
func Custom(options Options, statements ...Code) *Statement {
	return newStatement().Custom(options, statements...)
}

func (g *Group) Custom(options Options, statements ...Code) *Statement {
	return g.item(Custom(options, statements...))
}

func (s *Statement) Custom(options Options, statements ...Code) *Statement {
	return s.addGroup("custom", groupOp{
		open:  options.Open,
		close: options.Close,
		sep:   options.Separator,
		multi: options.Multi,
	}, statements)
}

// CustomFunc renders a customized statement list. Pass in options to specify
// multi-line, and tokens for open, close, separator.
func CustomFunc(options Options, f func(*Group)) *Statement {
	return newStatement().CustomFunc(options, f)
}

func (g *Group) CustomFunc(options Options, f func(*Group)) *Statement {
	return g.item(CustomFunc(options, f))
}

func (s *Statement) CustomFunc(options Options, f func(*Group)) *Statement {
	return s.addGroupFunc("custom", groupOp{
		open:  options.Open,
		close: options.Close,
		sep:   options.Separator,
		multi: options.Multi,
	}, f)
}
