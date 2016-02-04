package egon

import (
	"bytes"
	"fmt"
)

// TextBlock represents a UTF-8 encoded block of text that is written to the writer as-is.
type TextBlock struct {
	Pos     Pos
	Content string
}

func (b *TextBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)
	fmt.Fprintf(buf, `io.WriteString(w, %q)`+"\n", b.Content)
	return nil
}
