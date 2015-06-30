package egon_test

import (
	. "github.com/commondream/egon"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensure that a template can be written to a writer.
func TestTemplate_Write(t *testing.T) {
	tmpl := &Template{
		Path: "tmp/foo.egon",
		Blocks: []Block{
			&TextBlock{Content: "<html>", Pos: Pos{Path: "foo.ego", LineNo: 4}},
			&HeaderBlock{Content: "import \"fmt\"", Pos: Pos{Path: "tmp/foo.ego", LineNo: 8}},
			&ParameterBlock{ParamName: "nums", ParamType: "[]int"},
			&CodeBlock{Content: "  for _, num := range nums {"},
			&TextBlock{Content: "    <p>"},
			&RawPrintBlock{Content: "num + 1"},
			&TextBlock{Content: "    </p>"},
			&CodeBlock{Content: "  }"},
			&TextBlock{Content: "</html>"},
		},
	}
	p := &Package{Template: tmpl}
	err := p.Write()
	assert.NoError(t, err)
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
	_, err := tmpl.PackageName()
	assert.Error(t, err)
}

func TestTemplate_SourceFile(t *testing.T) {
	tmpl := &Template{Path: "foo.egon"}
	name := tmpl.SourceFile()

	assert.Equal(t, "foo.egon.go", name)
}
