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

func stripWhitespace(s string) string {
	out := make([]byte, 0)
	seenSpace := false
	for _, c := range s {
		switch c {
		case '\t', '\n', '\v', '\f', '\r':
				seenSpace = true
		case ' ':
			if !seenSpace {
				out = append(out, byte(c))
				seenSpace = true
			}
		default:
			out = append(out, byte(c))
			seenSpace = false
		}
	}
	return string(out)
}

func (b *TextBlock) write(buf *bytes.Buffer) error {
	if (Config.Minify) {
		b.Content = stripWhitespace(b.Content)
	}
	if (len(b.Content) > 0) {
		b.Pos.write(buf)
		fmt.Fprintf(buf, `io.WriteString(w, %q)`+"\n", b.Content)
	}
	return nil
}
