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
	TemplateName    = "/elm/Rpc.elm.gotmpl"
	ElmModuleOption = "elm_module"
)

func Generate(ns *spec.Namespace, outdir string) error {
	module := newModule(nil, outdir, ns)
	return writeModule(module)
}

func writeModule(mod *Module) error {
	if err := writeTmpl(mod.OutPath, TemplateName, mod); err != nil {
		return errors.Wrap(err, "elm template failure")
	}

	for _, child := range mod.Children {
		if err := writeModule(child); err != nil {
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
