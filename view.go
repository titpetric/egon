package egon

import (
	"io"
)

// View represents a runnable form of a template that can be passed between
// layers in an application without needing to render until the last necessary
// moment.
type View struct {
	PackageName  string
	Name         string
	TemplatePath string
	RenderFunc   func(io.Writer) error
}

// Render renders the view to the given io.Writer.
func (view *View) Render(w io.Writer) error {
	return view.RenderFunc(w)
}
