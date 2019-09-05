package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chakrit/rpc/generator"
	"github.com/chakrit/rpc/internal"
	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/parser"
	"github.com/chakrit/rpc/spec"
)

func main() {
	opts := parseOptions()
	logger := internal.NewLogger(opts.Silent)
	if err := opts.validate(); err != nil {
		logger.Fatal(err)
	}

	switch {
	case opts.LexOnly:
		lexMain(opts, logger)
	case opts.ParseOnly:
		parseMain(opts, logger)
	default:
		genMain(opts, logger)
	}
}

func lexMain(opts Options, logger internal.Logger) {
	var allTokens []*lexer.Token
	process(opts.SpecFilenames, logger, func(reader io.Reader) error {
		tokens, err := lexer.Lex(lexer.Options{
			Input:  reader,
			Logger: logger,
		})

		if err != nil {
			return err
		}

		allTokens = append(allTokens, tokens...)
		return nil
	})

	encoder := json.NewEncoder(os.Stdout)
	for _, token := range allTokens {
		if err := encoder.Encode(token); err != nil {
			logger.Fatal(fmt.Errorf("json encode failure: %w", err))
		}
	}
}

func parseMain(opts Options, logger internal.Logger) {
	root := &spec.Namespace{}
	process(opts.SpecFilenames, logger, func(reader io.Reader) error {
		if ns, err := parser.Parse(parser.Options{
			Input:  reader,
			Logger: logger,
		}); err != nil {
			return err
		} else {
			root = root.Merge(ns).(*spec.Namespace)
			return nil
		}
	})

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(root); err != nil {
		logger.Fatal(fmt.Errorf("json encode failure: %w", err))
	}
}

func genMain(opts Options, logger internal.Logger) {
	root := &spec.Namespace{}
	process(opts.SpecFilenames, logger, func(reader io.Reader) error {
		if ns, err := parser.Parse(parser.Options{
			Input:  reader,
			Logger: logger,
		}); err != nil {
			return err
		} else {
			root = root.Merge(ns).(*spec.Namespace)
			return nil
		}
	})

	err := generator.Generate(root, &generator.Options{
		Logger: logger,
		OutDir: opts.OutputDir,
		Target: opts.Target,
	})
	if err != nil {
		logger.Fatal(err)
	}
}

func process(filenames []string, logger internal.Logger, action func(io.Reader) error) {
	processOne := func(filename string) error { // provide scope for defer file closing
		file, err := openInput(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		return action(file)
	}

	for _, filename := range filenames {
		// fatal here, since we allow overrides from chaining multiple RPCs,
		// skipping one or more input file will have unintended side effects
		// so failing fast is the better option
		if err := processOne(filename); err != nil {
			logger.Fatal(fmt.Errorf("%s: %w", filename, err))
		}
	}
}

func openInput(filename string) (io.ReadCloser, error) {
	if filename == "" || filename == "-" {
		return os.Stdin, nil
	} else {
		return os.Open(filename)
	}
}
