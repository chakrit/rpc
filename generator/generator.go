package generator

//go:generate statik -src=./tmpl -dest=./ -p tmpldata -f

import (
	"errors"
	"fmt"

	"github.com/chakrit/rpc/generator/elm"
	"github.com/chakrit/rpc/generator/golang"
	"github.com/chakrit/rpc/internal"
	"github.com/chakrit/rpc/spec"
)

type Func func(ns *spec.Namespace, outdir string) error

type Interface interface {
	Generate(ns *spec.Namespace, outdir string) error
}

type Options struct {
	Logger internal.Logger
	OutDir string
	Target string
}

// added inside each implementation's init()
var implementations = map[string]Func{
	"elm": elm.Generate,
	"go":  golang.Generate,
}

func Generate(ns *spec.Namespace, opt *Options) error {
	generate, ok := implementations[opt.Target]
	if !ok {
		return errors.New("unsupported target `" + opt.Target + "`")
	}

	if err := generate(ns, opt.OutDir); err != nil {
		return fmt.Errorf("generator failure: %w", err)
	} else {
		return nil
	}
}
