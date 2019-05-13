package golang

import (
	"context"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/chakrit/rpc/generator/tmpldata"
	"github.com/chakrit/rpc/spec"
	"github.com/pkg/errors"
)

const (
	ImportOption  = "go_import"
	PackageOption = "go_package"
	OutName       = "rpc.go"

	SharedTemplateName = "/golang/shared.go.gotmpl"
	PkgTemplateName    = "/golang/pkg.go.gotmpl"
	ClientTemplateName = "/golang/client.go.gotmpl"
	ServerTemplateName = "/golang/server.go.gotmpl"

	DefaultPkgName    = "rpc"
	DefaultImportPath = "go.example.com/rpc"
)

func Generate(ns *spec.Namespace, outdir string) error {
	pkg := newRootPkg(ns)
	if err := writeRPCPackages(outdir, pkg); err != nil {
		return errors.Wrap(err, "go template failure")
	}
	if err := writeClientPackage(outdir, pkg); err != nil {
		return errors.Wrap(err, "go template failure")
	}
	if err := writeServerPackage(outdir, pkg); err != nil {
		return errors.Wrap(err, "go template failure")
	}

	return nil
}

func writeClientPackage(rootdir string, pkg *Pkg) error {
	return write(
		path.Join(rootdir, "client/client.go"),
		ClientTemplateName,
		pkg,
	)
}

func writeServerPackage(rootdir string, pkg *Pkg) error {
	return write(
		path.Join(rootdir, "server/server.go"),
		ServerTemplateName,
		pkg,
	)
}

func writeRPCPackages(rootdir string, pkg *Pkg) error {
	outpath := path.Join(rootdir, pkg.FilePath)
	if err := write(outpath, PkgTemplateName, pkg); err != nil {
		name := pkg.Name
		if name == "" {
			name = "root"
		}

		return errors.Wrap(err, "generating `"+outpath+"`")
	}

	for _, child := range pkg.Children {
		if err := writeRPCPackages(rootdir, child); err != nil {
			return err
		}
	}

	return nil
}

func write(outpath, tmplname string, pkg *Pkg) error {
	if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
		return errors.Wrap(err, "mkdir -p")
	}
	outfile, err := os.Create(outpath)
	if err != nil {
		return errors.Wrap(err, "creating `"+outpath+"`")
	}
	defer outfile.Close()

	tmplContent, err := tmpldata.Read(tmplname)
	if err != nil {
		return errors.Wrap(err, "reading template `"+tmplname+"`")
	}
	gotmpl, err := template.New(tmplname).Funcs(funcMap(pkg)).Parse(tmplContent)
	if err != nil {
		return errors.Wrap(err, "template parse")
	}

	typeRefContent, err := tmpldata.Read(SharedTemplateName)
	if err != nil {
		return errors.Wrap(err, "reading template `"+SharedTemplateName+"`")
	}
	gotmpl, err = gotmpl.Parse(typeRefContent)
	if err != nil {
		return errors.Wrap(err, "template parse")
	}

	err = gotmpl.Execute(outfile, pkg)
	if err != nil {
		return errors.Wrap(err, "template render")
	}

	if err := gofmt(outpath); err != nil {
		return errors.Wrap(err, "gofmt")
	}
	return nil
}

// TODO: A more generic means to run tools. Probably should not run `goimports` though as
//   it may get confused about the imports in the output folder and remove some lines
//   from the output code making it difficult to debug.
func gofmt(outpath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gofmt", "-s", "-w", outpath)
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
