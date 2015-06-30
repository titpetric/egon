package egon

import (
	"bytes"
)

// Block represents an element of the template.
type Block interface {
	write(*bytes.Buffer) error
}

// isTextBlock returns true if the block is a text block.
func isTextBlock(b Block) bool {
	_, ok := b.(*TextBlock)
	return ok
}
