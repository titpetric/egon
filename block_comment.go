package egon

import (
	"bytes"
)

// CommentBlock represents a block of text which is discarded
type CommentBlock struct {
	Pos     Pos
	Content string
}

func (b *CommentBlock) write(buf *bytes.Buffer) error {
	return nil
}
