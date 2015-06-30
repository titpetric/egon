package egon

import (
	"bytes"
	"fmt"
)

// PrintBlock represents a block that will HTML escape the contents before outputting
type PrintBlock struct {
	Pos     Pos
	Content string
}

func (b *PrintBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)
	fmt.Fprintf(buf, `_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%%v", %s)))`+"\n", b.Content)
	return nil
}
