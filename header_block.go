package egon

import (
	"bytes"
	"fmt"
)

// HeaderBlock represents a Go code block that is printed at the top of the template.
type HeaderBlock struct {
	Pos     Pos
	Content string
}

func (b *HeaderBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)
	fmt.Fprintln(buf, b.Content)
	return nil
}
