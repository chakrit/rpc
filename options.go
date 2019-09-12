package main

import (
	"flag"
	"os"
	"strings"

	"errors"
)

type Options struct {
	LexOnly   bool
	ParseOnly bool
	Silent    bool
	Target    string

	OutputDir     string
	CleanOutput   bool
	SpecFilenames []string
}

var (
	ErrNoInput      = errors.New("no rpc input file given")
	ErrNoGenTarget  = errors.New("no target specified for the generator")
	ErrNoOutput     = errors.New("no output folder specified for the generator")
	ErrOutputNotDir = errors.New("specified output path is not a directory")
)

func parseOptions() Options {
	var options = Options{}
	flag.BoolVar(&options.Silent, "q", false, "Silence all warnings.")
	flag.BoolVar(&options.LexOnly, "lex", false, "Lex MRPC file and print a list of tokens found.")
	flag.BoolVar(&options.ParseOnly, "parse", false, "Parse MRPC file and output a JSON spec for further processing.")
	flag.StringVar(&options.Target, "gen", "", "Generate an implementation for the specified target.")
	flag.BoolVar(&options.CleanOutput, "clean", false, "Cleans the output directory before generating new files into it.")
	flag.StringVar(&options.OutputDir, "out", "", "Output directory or filename. Defaults to STDOUT.")
	flag.Parse()

	options.OutputDir = strings.TrimSpace(options.OutputDir)
	options.Target = strings.TrimSpace(options.Target)
	for _, arg := range flag.Args() {
		options.SpecFilenames = append(options.SpecFilenames, normalizeFilename(arg))
	}

	return options
}

func (opts Options) validate() error {
	genMode := !opts.ParseOnly && !opts.LexOnly

	switch {
	case len(opts.SpecFilenames) == 0:
		return ErrNoInput
	case genMode && opts.Target == "":
		return ErrNoGenTarget
	case genMode && opts.OutputDir == "":
		return ErrNoOutput
	}

	if info, err := os.Stat(opts.OutputDir); err != nil {
		return err
	} else if !info.IsDir() {
		return ErrOutputNotDir
	}

	return nil
}

func normalizeFilename(str string) string {
	str = strings.TrimSpace(str)
	if str == "-" { // we default to STDIN STDOUT already, so "-" should have no effect
		str = ""
	}
	return str
}
