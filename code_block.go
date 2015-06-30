package egon

import (
	"bytes"
	"fmt"
)

// CodeBlock represents a Go code block that is printed as-is to the template.
type CodeBlock struct {
	Pos     Pos
	Content string
}

func (b *CodeBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)
	fmt.Fprintln(buf, b.Content)
	return nil
}
