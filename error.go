package egon

import "errors"

var (
	// ErrParameterFormat notifies the user that a parameter tag is poorly formatted.
	ErrParameterFormat = errors.New("parameters should be of form `param type`")
)
