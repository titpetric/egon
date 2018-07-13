package egon

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

// Template represents an entire Ego template.
// Templates consist of a set of parameters and other block.
// Blocks can be either a TextBlock, a PrintBlock, a RawPrintBlock, or a CodeBlock.
type Template struct {
	Path   string
	Blocks []Block
}

// PackageName returns the name of the package, based on the last non-file
// part of the path.
func (t *Template) PackageName() (string, error) {
	path, err := filepath.Abs(t.Path)
	if err != nil {
		return "", ErrUnidentifiablePackage
	}

	// split the path by file separator, rip the first one off (it's always blank)
	// and then grab the last one
	parts := strings.Split(path, string(filepath.Separator))
	parts = parts[1:]
	if len(parts) < 2 {
		return "", ErrUnidentifiablePackage
	}
	return parts[len(parts)-2], nil
}

// FileName returns the filename of the template, without the path.
func (t *Template) FileName() string {
	_, fileName := filepath.Split(t.Path)
	return fileName
}

// Name returns a name for the template as a camel cased string based on
// the filename.
func (t *Template) Name() string {
	fileName := t.FileName()

	// remove the extension
	parts := strings.Split(fileName, ".")
	name := parts[0]

	// Filter out any non-letter and digit runes
	re := regexp.MustCompile("[^\\p{L}0-9]")
	name = re.ReplaceAllString(name, " ")

	// convert to title case and remove spaces
	name = strings.Title(name)
	name = strings.Replace(name, " ", "", -1)

	return name
}

// TemplateFuncName returns the name of the Template func for this template.
func (t *Template) TemplateFuncName() string {
	return strings.Join([]string{t.Name(), "Template"}, "")
}

// ViewFuncName returns the name fo the View func for this template.
func (t *Template) ViewFuncName() string {
	return strings.Join([]string{t.Name(), "View"}, "")
}

// SourceFile returns the path to the source file that should be
// generated from this template.
func (t *Template) SourceFile() string {
	return strings.Join([]string{t.Path, ".go"}, "")
}

// Write writes the template to a writer.
func (t *Template) Write(w io.Writer) error {
	buf := new(bytes.Buffer)

	if err := t.writeHeader(buf); err != nil {
		return err
	}

	params := t.parameterBlocks()
	buf.WriteString("\n")

	// render the template func
	// add the writer param
	ioParam := ParameterBlock{ParamName: "w", ParamType: "io.Writer"}
	params = append([]*ParameterBlock{&ioParam}, params...)
	buf.WriteString(fmt.Sprintf("func %s(", t.TemplateFuncName()))
	t.writeParameters(buf, params)
	buf.WriteString(") error {\n")

	// Write non-header blocks.
	for _, b := range t.nonHeaderBlocks() {
		if err := b.write(buf); err != nil {
			return err
		}
	}

	// Write return and function closing brace.
	buf.WriteString("return nil\n")
	buf.WriteString("}\n\n")

	// Write a simple `fn() string` function
	/*
		fn := t.TemplateFuncName() + "String"
		buf.WriteString(fmt.Sprintf("func %s() string {\n", fn))
		buf.WriteString("\tbuf = new(bytes.Buffer)\n")
		buf.WriteString("\t" + fmt.Sprintf("%s(buf)\n", fn))
		buf.WriteString("\treturn buf.String()\n")
		buf.WriteString("}\n")
	*/

	// Write code to external writer.
	_, err := buf.WriteTo(w)
	return err
}

func (t *Template) String() string {
	buf := new(bytes.Buffer)
	t.Write(buf)
	return buf.String()
}

func (t *Template) writeParameters(buf *bytes.Buffer, params []*ParameterBlock) {
	maxIndex := len(params) - 1
	for i, param := range params {
		param.write(buf)

		if i < maxIndex {
			buf.WriteString(", ")
		}
	}
}

// Writes the package name and consolidated header blocks.
func (t *Template) writeHeader(w io.Writer) error {
	name, err := t.PackageName()
	if err != nil {
		return err
	}

	// Write naive header first.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "package %s\n", name)
	for _, b := range t.headerBlocks() {
		b.write(&buf)
	}

	// Parse header into Go AST.
	f, err := parser.ParseFile(token.NewFileSet(), "ego.go", buf.String(), parser.ImportsOnly)
	if err != nil {
		fmt.Println(buf.String())
		return fmt.Errorf("writeHeader: %s", err)
	}

	// Reset buffer.
	buf.Reset()

	// Write deduped imports.
	var decls = map[string]bool{`:"fmt"`: true, `:"io"`: true}
	fmt.Fprint(&buf, "import (\n")
	if t.hasFmtPrintBlock() {
		fmt.Fprintln(&buf, `"fmt"`)
	}
	if t.hasItoaPrintBlock() {
		fmt.Fprintln(&buf, `"strconv"`)
	}
	if t.hasEscapedPrintBlock() {
		fmt.Fprintln(&buf, `"html"`)
		decls["html"] = true
	}
	fmt.Fprintln(&buf, `"io"`)

	for _, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			continue
		}

		for _, s := range d.Specs {
			s := s.(*ast.ImportSpec)
			var id string
			if s.Name != nil {
				id = s.Name.Name
			}
			id += ":" + s.Path.Value

			// Ignore any imports which have already been imported.
			if decls[id] {
				continue
			}
			decls[id] = true

			// Otherwise write it.
			if s.Name == nil {
				fmt.Fprintf(&buf, "%s\n", s.Path.Value)
			} else {
				fmt.Fprintf(&buf, "%s %s\n", s.Name.Name, s.Path.Value)
			}
		}
	}
	fmt.Fprint(&buf, ")\n")

	// Write out to writer.
	buf.WriteTo(w)

	return nil
}

func (t *Template) parameterBlocks() []*ParameterBlock {
	blocks := []*ParameterBlock{}
	for _, b := range t.Blocks {
		if b, ok := b.(*ParameterBlock); ok {
			blocks = append(blocks, b)
		}
	}
	return blocks
}

func (t *Template) headerBlocks() []*HeaderBlock {
	var blocks []*HeaderBlock
	for _, b := range t.Blocks {
		if b, ok := b.(*HeaderBlock); ok {
			blocks = append(blocks, b)
		}
	}
	return blocks
}

func (t *Template) nonHeaderBlocks() []Block {
	var blocks []Block
	for _, b := range t.Blocks {
		switch b.(type) {
		case *ParameterBlock, *HeaderBlock:
		default:
			blocks = append(blocks, b)
		}
	}
	return blocks
}

func (t *Template) hasEscapedPrintBlock() bool {
	for _, b := range t.Blocks {
		if _, ok := b.(*PrintBlock); ok {

			return true
		}
	}
	return false
}

func (t *Template) hasItoaPrintBlock() bool {
	for _, b := range t.Blocks {
		if pBlock, ok := b.(*PrintBlock); ok {
			if pBlock.Type == 'd' {
				return true
			}
		}
	}
	return false
}

func (t *Template) hasFmtPrintBlock() bool {
	for _, b := range t.Blocks {
		if pBlock, ok := b.(*PrintBlock); ok {
			if pBlock.Type != 'd' && pBlock.Type != 's' {
				return true
			}
		}
	}
	return false
}

// normalize joins together adjacent text blocks.
func (t *Template) normalize() {
	var a []Block
	for _, b := range t.Blocks {
		if isTextBlock(b) && len(a) > 0 && isTextBlock(a[len(a)-1]) {
			a[len(a)-1].(*TextBlock).Content += b.(*TextBlock).Content
		} else {
			a = append(a, b)
		}
	}
	t.Blocks = a
}

var _ fmt.Stringer = &Template{}
