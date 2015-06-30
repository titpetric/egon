package egon_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	. "github.com/commondream/egon"
	"github.com/stretchr/testify/assert"
)

// Ensure that a template can be written to a writer.
func TestTemplate_Write(t *testing.T) {
	var buf bytes.Buffer
	tmpl := &Template{
		Path: "/some/path/to/foo.egon",
		Blocks: []Block{
			&TextBlock{Content: "<html>", Pos: Pos{Path: "foo.ego", LineNo: 4}},
			&HeaderBlock{Content: "import \"fmt\"", Pos: Pos{Path: "foo.ego", LineNo: 8}},
			&DeclarationBlock{ParamName: "nums", ParamType: "[]int"},
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
	buf.WriteTo(os.Stdout)
}

func warn(v ...interface{})              { fmt.Fprintln(os.Stderr, v...) }
func warnf(msg string, v ...interface{}) { fmt.Fprintf(os.Stderr, msg+"\n", v...) }
