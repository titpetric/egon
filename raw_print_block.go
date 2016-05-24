package egon

import (
	"bytes"
	"fmt"
)

// RawPrintBlock represents a block of the template that is printed out to the writer.
type RawPrintBlock struct {
	Pos     Pos
	Content string
	Type    byte
}

func (b *RawPrintBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)
	fmt.Fprintf(buf, `io.WriteString(w, %s)`+"\n", b.Content)
	return nil
}
