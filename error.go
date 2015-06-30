package ego

import "errors"

var (
	// ErrDeclarationRequired is returned when there is no declaration block
	// in a template.
	ErrDeclarationRequired = errors.New("declaration required")
	ErrDeclarationFormat   = errors.New("declarations should be of form `param type`")
)
