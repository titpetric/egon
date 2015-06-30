package egon_test

import (
	"bytes"
	. "github.com/commondream/egon"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// Ensure that a template can be written to a writer.
func TestTemplate_Write(t *testing.T) {
	var buf bytes.Buffer
	tmpl := &Template{
		Path: "/some/path/to/foo.egon",
		Blocks: []Block{
			&TextBlock{Content: "<html>", Pos: Pos{Path: "foo.ego", LineNo: 4}},
			&HeaderBlock{Content: "import \"fmt\"", Pos: Pos{Path: "foo.ego", LineNo: 8}},
			&ParameterBlock{ParamName: "nums", ParamType: "[]int"},
			&CodeBlock{Content: "  for _, num := range nums {"},
			&TextBlock{Content: "    <p>"},
			&RawPrintBlock{Content: "num + 1"},
			&TextBlock{Content: "    </p>"},
			&CodeBlock{Content: "  }"},
			&TextBlock{Content: "</html>"},
		},
	}
	p := &Package{Templates: []*Template{tmpl}, Name: "foo"}
	err := p.Write(&buf)
	assert.NoError(t, err)
	//buf.WriteTo(os.Stdout)
}

// PackageName Tests
func TestTemplate_PackageName(t *testing.T) {
	tmpl := &Template{Path: "/some/path/to/foo.egon"}
	name, err := tmpl.PackageName()

	assert.NoError(t, err)
	assert.Equal(t, "to", name)
}

func TestTemplate_PackageNameRelative(t *testing.T) {
	tmpl := &Template{Path: "template.go"}
	name, err := tmpl.PackageName()

	assert.NoError(t, err)
	assert.Equal(t, "egon", name)
}

func TestTemplate_PackageNameNoFolder(t *testing.T) {
	tmpl := &Template{Path: "/foo.egon"}
	name, err := tmpl.PackageName()
	log.Println(name)
	assert.Error(t, err)
}
