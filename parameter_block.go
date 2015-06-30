package egon

import (
	"bytes"
	"fmt"
)

// DeclarationBlock represents a block that declaration the function signature.
type ParameterBlock struct {
	Pos       Pos
	ParamName string
	ParamType string
}

func (b *ParameterBlock) write(buf *bytes.Buffer) error {
	b.Pos.write(buf)
	fmt.Fprintf(buf, "%s %s", b.ParamName, b.ParamType)
	return nil
}
