package egon

import (
	"bytes"
	"fmt"
)

// Pos represents a position in a given file.
type Pos struct {
	Path   string
	LineNo int
}

func (p *Pos) write(buf *bytes.Buffer) {
	if p != nil && p.Path != "" && p.LineNo > 0 {
		fmt.Fprintf(buf, "//line %s:%d\n", p.Path, p.LineNo)
	}
}
