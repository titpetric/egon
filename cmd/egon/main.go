package main

import (
	"flag"
	"fmt"
	"go/scanner"
	"log"
	"os"
	"path/filepath"

	"github.com/commondream/egon"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	// If no paths are provided then use the present working directory.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Recursively retrieve all ego templates
	var v visitor
	for _, root := range roots {
		if err := filepath.Walk(root, v.visit); err != nil {
			scanner.PrintError(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Parse every *.ego file.
	for _, path := range v.paths {
		template, err := egon.ParseFile(path)
		if err != nil {
			log.Fatal("parse file: ", err)
		}

		pkg := &egon.Package{Template: template}
		err = pkg.Write()
		if err != nil {
			log.Fatal("write: ", err)
		}
	}
}

// visitor iterates over
type visitor struct {
	paths []string
}

func (v *visitor) visit(path string, info os.FileInfo, err error) error {
	if info == nil {
		return fmt.Errorf("file not found: %s", path)
	}
	if !info.IsDir() && filepath.Ext(path) == ".egon" {
		v.paths = append(v.paths, path)
	}
	return nil
}
