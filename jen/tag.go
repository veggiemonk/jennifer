package jen

import (
	"fmt"
	"io"
	"slices"
	"strconv"
)

// Tag renders a struct tag.
func Tag(items map[string]string) *Statement            { return newStatement().Tag(items) }
func (g *Group) Tag(items map[string]string) *Statement { return g.item(Tag(items)) }
func (s *Statement) Tag(items map[string]string) *Statement {
	*s = append(*s, tag{items: items})
	return s
}

type tag struct {
	items map[string]string
}

func (t tag) isNull(f *File) bool {
	return len(t.items) == 0
}

func (t tag) render(f *File, w io.Writer, s *Statement) error {
	if t.isNull(f) {
		return nil
	}

	sorted := make([]string, 0, len(t.items))
	for k := range t.items {
		sorted = append(sorted, k)
	}
	slices.Sort(sorted)

	var str string
	for _, k := range sorted {
		v := t.items[k]
		if len(str) > 0 {
			str += " "
		}
		str += fmt.Sprintf(`%s:%q`, k, v)
	}

	if strconv.CanBackquote(str) {
		str = "`" + str + "`"
	} else {
		str = strconv.Quote(str)
	}

	_, err := w.Write([]byte(str))
	return err
}
