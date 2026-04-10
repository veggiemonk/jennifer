package jen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var nonAlphanumRegex = regexp.MustCompile(`[^a-z0-9]`)

// NewFile Creates a new file, with the specified package name.
func NewFile(packageName string) *File {
	return &File{
		Group:   &Group{multi: true},
		name:    packageName,
		imports: map[string]importdef{},
		hints:   map[string]importdef{},
	}
}

// NewFilePath creates a new file while specifying the package path - the
// package name is inferred from the path.
func NewFilePath(packagePath string) *File {
	return &File{
		Group:   &Group{multi: true},
		name:    guessAlias(packagePath),
		path:    packagePath,
		imports: map[string]importdef{},
		hints:   map[string]importdef{},
	}
}

// NewFilePathName creates a new file with the specified package path and name.
func NewFilePathName(packagePath, packageName string) *File {
	return &File{
		Group:   &Group{multi: true},
		name:    packageName,
		path:    packagePath,
		imports: map[string]importdef{},
		hints:   map[string]importdef{},
	}
}

// File represents a single source file. Package imports are managed
// automatically by File.
type File struct {
	*Group
	name        string
	path        string
	imports     map[string]importdef
	hints       map[string]importdef
	comments    []string
	headers     []string
	cgoPreamble []string
	// NoFormat can be set to true to disable formatting of the generated source. This may be useful
	// when performance is critical, and readable code is not required.
	NoFormat bool
	// If you're worried about generated package aliases conflicting with local variable names, you
	// can set a prefix here. Package foo becomes {prefix}_foo.
	PackagePrefix string
	// CanonicalPath adds a canonical import path annotation to the package clause.
	CanonicalPath string
}

// importdef differentiates packages where we know the name from packages where
// the import is aliased.
type importdef struct {
	name  string
	alias bool
}

// HeaderComment adds a comment to the top of the file, above any package
// comments. A blank line is rendered below the header comments, ensuring
// header comments are not included in the package doc.
func (f *File) HeaderComment(comment string) {
	f.headers = append(f.headers, comment)
}

// PackageComment adds a comment to the top of the file, above the package
// keyword.
func (f *File) PackageComment(comment string) {
	f.comments = append(f.comments, comment)
}

// CgoPreamble adds a cgo preamble comment that is rendered directly before the
// "C" pseudo-package import.
func (f *File) CgoPreamble(comment string) {
	f.cgoPreamble = append(f.cgoPreamble, comment)
}

// Anon adds an anonymous import.
func (f *File) Anon(paths ...string) {
	for _, p := range paths {
		f.imports[p] = importdef{name: "_", alias: true}
	}
}

// ImportName provides the package name for a path. If specified, the alias will
// be omitted from the import block.
func (f *File) ImportName(path, name string) {
	f.hints[path] = importdef{name: name, alias: false}
}

// ImportNames allows multiple names to be imported as a map.
func (f *File) ImportNames(names map[string]string) {
	for path, name := range names {
		f.hints[path] = importdef{name: name, alias: false}
	}
}

// ImportAlias provides the alias for a package path that should be used in the
// import block. A period can be used to force a dot-import.
func (f *File) ImportAlias(path, alias string) {
	f.hints[path] = importdef{name: alias, alias: true}
}

func (f *File) isLocal(path string) bool {
	return f.path == path
}

func (f *File) isValidAlias(alias string) bool {
	if alias == "." {
		return true
	}
	if IsReservedWord(alias) {
		return false
	}
	for _, v := range f.imports {
		if alias == v.name {
			return false
		}
	}
	return true
}

func (f *File) isDotImport(path string) bool {
	if id, ok := f.hints[path]; ok {
		return id.name == "." && id.alias
	}
	return false
}

func (f *File) register(path string) string {
	if f.isLocal(path) {
		return ""
	}

	// already registered
	def := f.imports[path]
	if def.name != "" && def.name != "_" {
		return def.name
	}

	// "C" pseudo-package
	if path == "C" {
		f.imports["C"] = importdef{name: "C", alias: false}
		return "C"
	}

	var name string
	var alias bool

	if hint := f.hints[path]; hint.name != "" {
		name = hint.name
		alias = hint.alias
	} else if standardLibraryHints[path] != "" {
		name = standardLibraryHints[path]
		alias = false
	} else {
		name = guessAlias(path)
		alias = true
	}

	// make unique
	unique := name
	i := 0
	for !f.isValidAlias(unique) {
		i++
		unique = fmt.Sprintf("%s%d", name, i)
		alias = true
	}

	if f.PackagePrefix != "" && alias {
		unique = f.PackagePrefix + "_" + unique
	}

	f.imports[path] = importdef{name: unique, alias: alias}
	return unique
}

// Save renders the file and saves to the filename provided.
func (f *File) Save(filename string) error {
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return fmt.Errorf("saving %s: %w", filename, err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", filename, err)
	}
	return nil
}

// Render renders the file to the provided writer.
func (f *File) Render(w io.Writer) error {
	body := &bytes.Buffer{}
	if err := f.render(f, body, nil); err != nil {
		return fmt.Errorf("rendering body: %w", err)
	}
	source := &bytes.Buffer{}
	if len(f.headers) > 0 {
		for _, c := range f.headers {
			if err := Comment(c).render(f, source, nil); err != nil {
				return err
			}
			if _, err := fmt.Fprint(source, "\n"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(source, "\n"); err != nil {
			return err
		}
	}
	for _, c := range f.comments {
		if err := Comment(c).render(f, source, nil); err != nil {
			return err
		}
		if _, err := fmt.Fprint(source, "\n"); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(source, "package %s", f.name); err != nil {
		return err
	}
	if f.CanonicalPath != "" {
		if _, err := fmt.Fprintf(source, " // import %q", f.CanonicalPath); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(source, "\n\n"); err != nil {
		return err
	}
	if err := f.renderImports(source); err != nil {
		return err
	}
	if _, err := source.Write(body.Bytes()); err != nil {
		return err
	}
	var output []byte
	if f.NoFormat {
		output = source.Bytes()
	} else {
		var err error
		output, err = format.Source(source.Bytes())
		if err != nil {
			return fmt.Errorf("formatting generated source: %w\n%s", err, numberLines(source.String()))
		}
	}
	if _, err := w.Write(output); err != nil {
		return err
	}
	return nil
}

func (f *File) renderImports(source io.Writer) error {
	hasCgo := f.imports["C"].name != "" || len(f.cgoPreamble) > 0
	separateCgo := hasCgo && len(f.cgoPreamble) > 0

	filtered := map[string]importdef{}
	for path, def := range f.imports {
		if path == "C" && separateCgo {
			continue
		}
		filtered[path] = def
	}

	if len(filtered) == 1 {
		for path, def := range filtered {
			if def.alias && path != "C" {
				if _, err := fmt.Fprintf(source, "import %s %s\n\n", def.name, strconv.Quote(path)); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(source, "import %s\n\n", strconv.Quote(path)); err != nil {
					return err
				}
			}
		}
	} else if len(filtered) > 1 {
		if _, err := fmt.Fprint(source, "import (\n"); err != nil {
			return err
		}
		paths := make([]string, 0, len(filtered))
		for path := range filtered {
			paths = append(paths, path)
		}
		slices.Sort(paths)
		for _, path := range paths {
			def := filtered[path]
			if def.alias && path != "C" {
				if _, err := fmt.Fprintf(source, "%s %s\n", def.name, strconv.Quote(path)); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(source, "%s\n", strconv.Quote(path)); err != nil {
					return err
				}
			}
		}
		if _, err := fmt.Fprint(source, ")\n\n"); err != nil {
			return err
		}
	}

	if separateCgo {
		for _, c := range f.cgoPreamble {
			if err := Comment(c).render(f, source, nil); err != nil {
				return err
			}
			if _, err := fmt.Fprint(source, "\n"); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(source, "import \"C\"\n\n"); err != nil {
			return err
		}
	}

	return nil
}

// GoString renders the File for testing. Any error will cause a panic.
func (f *File) GoString() string {
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		panic(fmt.Errorf("jennifer: File.GoString render error: %w", err))
	}
	return buf.String()
}

func guessAlias(path string) string {
	alias := strings.TrimSuffix(path, "/")
	if strings.Contains(alias, "/") {
		alias = alias[strings.LastIndex(alias, "/")+1:]
	}
	alias = strings.ToLower(alias)
	alias = nonAlphanumRegex.ReplaceAllString(alias, "")
	for firstRune, runeLen := utf8.DecodeRuneInString(alias); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(alias) {
		alias = alias[runeLen:]
	}
	if alias == "" {
		alias = "pkg"
	}
	return alias
}
