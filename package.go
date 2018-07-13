package egon

import (
	"fmt"
	"os"
)

// Package represents the source file representation of a template.
// Note that it's on its way out, in favor of this functionality being included
// in Template.
type Package struct {
	Template *Template
}

// Write writes out the package header and templates to a writer.
func (p *Package) Write() error {
	f, err := os.Create(p.Template.SourceFile())
	defer f.Close()

	if err != nil {
		return err
	}

	if err := p.Template.Write(f); err != nil {
		return fmt.Errorf("template: %s: %s", p.Template.Path, err)
	}

	return nil
}

