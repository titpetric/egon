package egon

import (
	"bytes"
	"fmt"
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
	name = strings.Join([]string{name, "Template"}, "")

	return name
}

// Write writes the template to a writer.
func (t *Template) Write(w io.Writer) error {
	var buf bytes.Buffer

	params := t.parameterBlocks()

	// add the writer param
	ioParam := ParameterBlock{ParamName: "w", ParamType: "io.Writer"}
	params = append([]*ParameterBlock{&ioParam}, params...)

	buf.WriteString(fmt.Sprintf("func %s(", t.Name()))
	maxIndex := len(params) - 1
	for i, param := range params {
		param.write(&buf)

		if i < maxIndex {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(") {")

	// Write non-header blocks.
	for _, b := range t.nonHeaderBlocks() {
		if err := b.write(&buf); err != nil {
			return err
		}
	}

	// Write return and function closing brace.
	fmt.Fprint(&buf, "return nil\n")
	fmt.Fprint(&buf, "}\n")

	// Write code to external writer.
	_, err := buf.WriteTo(w)
	return err
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
