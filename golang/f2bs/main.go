package main

import (
	"flag"
	"io"
	"os"
)

var (
	inputFileName  = flag.String("input", "", "input filename")
	outputFileName = flag.String("output", "", "output filename")
	packageName    = flag.String("package", "main", "package name")
	varName        = flag.String("var", "_", "variable name")
	compress       = flag.Bool("compress", false, "use gzip compression")
	buildTags      = flag.String("buildtags", "", "build tags")
)

func run() error {
	var out io.Writer
	if *outputFileName != "" {
		f, err := os.Create(*outputFileName)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	var in io.Reader
	if *inputFileName != "" {
		f, err := os.Open(*inputFileName)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	if err := Write(out, in, *compress, *buildTags, *packageName, *varName); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
