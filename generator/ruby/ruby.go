package ruby

import (
	"os"
	"path"
	"text/template"

	"github.com/chakrit/rpc/generator/tmpldata"
	"github.com/chakrit/rpc/spec"
	"github.com/pkg/errors"
)

const (
	TemplateName = "/ruby/impl.rb.gotmpl"
	OutName      = "rpc.rb"
)

func Generate(ns *spec.Namespace, outdir string) error {
	rbtmpl, err := loadTemplate()
	if err != nil {
		return errors.Wrap(err, "ruby template failure")
	}

	if err := os.MkdirAll(outdir, 0755); err != nil {
		return errors.Wrap(err, "failed to create `"+outdir+"`")
	}

	dest := path.Join(outdir, OutName)
	outfile, err := os.Create(dest)
	if err != nil {
		return errors.Wrap(err, "failed to create `"+dest+"`")
	}
	defer outfile.Close()

	return rbtmpl.Execute(outfile, ns)
}

func loadTemplate() (*template.Template, error) {
	content, err := tmpldata.Read(TemplateName)
	if err != nil {
		return nil, err
	}

	rbtmpl, err := template.New(TemplateName).Parse(content)
	if err != nil {
		return nil, err
	}

	return rbtmpl, nil
}
