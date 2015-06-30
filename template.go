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

// Name returns a name for the template as a filename.
func (t *Template) Name() string {
	_, fileName := filepath.Split(t.Path)
	parts := strings.Split(fileName, ".")
	name := parts[0]
	re := regexp.MustCompile("[^\\p{L}0-9]")
	name = re.ReplaceAllString(name, " ")
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
