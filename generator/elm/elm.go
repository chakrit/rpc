package elm

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/chakrit/rpc/generator/tmpldata"
	"github.com/chakrit/rpc/spec"
)

const (
	RpcTemplateName  = "/elm/Rpc.elm.gotmpl"
	UtilTemplateName = "/elm/RpcUtil.elm.gotmpl"
	ModuleOption     = "elm_module"
)

type (
	Field struct {
		Name string
		Type *TypeRef
	}

	Member struct {
		Name  string
		Value string
	}

	Type struct {
		Name   string
		Fields []*Field
		Module *Module
	}

	Enum struct {
		Name    string
		Members []*Member
		Module  *Module
	}

	TypeRef struct {
		Name   string
		Args   []*TypeRef
		Module *Module
	}

	TypeResolution struct {
		Name    string
		Encode  string
		Decode  string
		Default string
	}

	Tuple struct {
		Name string
		Args []*TypeRef
	}

	RpcFunc struct {
		Name    string
		RPCPath string

		InArgs  []*TypeRef
		OutArgs []*TypeRef
	}
)

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
		return fmt.Errorf("elm template failure: %w", err)
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
