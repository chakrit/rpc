package elm

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/chakrit/rpc/generator/tmpldata"

	"github.com/chakrit/rpc/spec"
	"github.com/pkg/errors"
)

const (
	RpcTemplateName  = "/elm/Rpc.elm.gotmpl"
	UtilTemplateName = "/elm/RpcUtil.elm.gotmpl"
	ModuleOption     = "elm_module"
)

type ElmField struct {
	Name string
	Type *ElmTypeRef
}

type ElmType struct {
	Name   string
	Fields []*ElmField
	Module *Module
}

type ElmTypeRef struct {
	Name   string
	Args   []*ElmTypeRef
	Module *Module
}

type ElmTypeResolution struct {
	Name   string
	Encode string
	Decode string
}

type ElmTuple struct {
	Name string
	Args []*ElmTypeRef
}

type ElmRPCFunc struct {
	Name    string
	RPCPath string

	InArgs  []*ElmTypeRef
	OutArgs []*ElmTypeRef
}

func Generate(ns *spec.Namespace, outdir string) error {
	module := newModule(nil, outdir, ns)
	utilModule := newUtilModule(module, outdir, ns)

	if err := writeModule(utilModule, UtilTemplateName); err != nil {
		return err
	}
	if err := writeModule(module, RpcTemplateName); err != nil {
		return err
	}
	return nil
}

func writeModule(mod *Module, templateName string) error {
	if err := writeTmpl(mod.OutPath, templateName, mod); err != nil {
		return errors.Wrap(err, "elm template failure")
	}

	for _, child := range mod.Children {
		if err := writeModule(child, templateName); err != nil {
			return err
		}
	}

	return nil
}

func writeTmpl(outpath, tmplname string, mod *Module) error {
	if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
		return err
	}

	outfile, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	tmplContent, err := tmpldata.Read(tmplname)
	if err != nil {
		return err
	}

	tmpl, err := template.New(tmplname).Funcs(funcMap()).Parse(tmplContent)
	if err != nil {
		return err
	}

	return tmpl.Execute(outfile, mod)
}
