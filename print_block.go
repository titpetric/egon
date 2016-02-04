package egon

import (
	"bytes"
	"fmt"
)

// PrintBlock represents a block that will HTML escape the contents before outputting
type PrintBlock struct {
	Pos      Pos
	Content  string
	Type     byte
}

func (b *PrintBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)

	switch b.Type {
	case 'd':
		fmt.Fprintf(buf, `io.WriteString(w, strconv.Itoa(%s))`+"\n", b.Content)
	case 's':
		fmt.Fprintf(buf, `io.WriteString(w, html.EscapeString(%s))`+"\n", b.Content)
	case 0:
		fmt.Fprintf(buf, `io.WriteString(w, html.EscapeString(fmt.Sprintf("%%v", %s)))`+"\n", b.Content)
	default:
		fmt.Fprintf(buf, `io.WriteString(w, html.EscapeString(fmt.Sprintf("%%%c", %s)))`+"\n", b.Type, b.Content)
	}
	return nil
}
