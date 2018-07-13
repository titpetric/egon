package main

import (
	"fmt"
	"go/scanner"
	"log"
	"os"
	"path/filepath"

	"github.com/titpetric/egon"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	kingpin.Version("0.9.0")
	kingpin.Flag("extension", "templatefile extension").Short('e').Default("egon").StringVar(&egon.Config.TmplExtension)
	kingpin.Flag("typesafe", "if present use provided flags for format-string").Short('t').Default("true").BoolVar(&egon.Config.Typesafe)
	kingpin.Flag("stropt", "optimise string handling to reduce allocations").Short('s').Default("true").BoolVar(&egon.Config.StringOptimisations)
	kingpin.Flag("debug", "include debug comments in generated code").Short('d').Default("false").BoolVar(&egon.Config.Debug)
	kingpin.Flag("minify", "remove whitespace from output").Short('m').Default("false").BoolVar(&egon.Config.Minify)
	kingpin.Arg("folders", "folders to be processed").StringsVar(&egon.Config.Folders)
}

func main() {
	log.SetFlags(0)
	kingpin.CommandLine.Help = "Generate native Go code from ERB-style Templates"
	kingpin.Parse()

	if len(egon.Config.Folders) == 0 {
		egon.Config.Folders = []string{"."}
	}

	// Recursively retrieve all templates
	var v visitor
	for _, root := range egon.Config.Folders {
		log.Printf("scanning folder [%s]", root)
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
	if !info.IsDir() && filepath.Ext(path) == ("."+egon.Config.TmplExtension) {
		v.paths = append(v.paths, path)
	}
	return nil
}
