package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/chakrit/rpc/generator/tmpldata"

	"github.com/chakrit/rpc/spec"
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
		return fmt.Errorf("go template failure: %w", err)
	}
	if err := writeClientPackage(outdir, pkg); err != nil {
		return fmt.Errorf("go template failure: %w", err)
	}
	if err := writeServerPackage(outdir, pkg); err != nil {
		return fmt.Errorf("go template failure: %w", err)
	}

	return nil
}

func writeClientPackage(rootdir string, pkg *Pkg) error {
	return write(
		path.Join(rootdir, "client/client.go"),
		ClientTemplateName,
		pkg.Registry,
		&PkgContext{
			ContextPkg: &Pkg{Name: "client", MangledName: "client", ImportPath: "client"},
			DataPkg:    pkg,
		},
	)
}

func writeServerPackage(rootdir string, pkg *Pkg) error {
	return write(
		path.Join(rootdir, "server/server.go"),
		ServerTemplateName,
		pkg.Registry,
		&PkgContext{
			ContextPkg: &Pkg{Name: "server", MangledName: "server", ImportPath: "server"},
			DataPkg:    pkg,
		},
	)
}

func writeRPCPackages(rootdir string, pkg *Pkg) error {
	outpath := path.Join(rootdir, pkg.FilePath)
	if err := write(outpath, PkgTemplateName, pkg.Registry, pkg); err != nil {
		name := pkg.Name
		if name == "" {
			name = "root"
		}

		return fmt.Errorf("generating `"+outpath+"`: %w", err)
	}

	for _, child := range pkg.Children {
		if err := writeRPCPackages(rootdir, child); err != nil {
			return err
		}
	}

	return nil
}

func write(outpath, tmplname string, registry TypeRegistry, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
		return fmt.Errorf("mkdir -p: %w", err)
	}
	outfile, err := os.Create(outpath)
	if err != nil {
		return fmt.Errorf("creating `"+outpath+"`: %w", err)
	}
	defer outfile.Close()

	tmplContent, err := tmpldata.Read(tmplname)
	if err != nil {
		return fmt.Errorf("reading template `"+tmplname+"`: %w", err)
	}
	gotmpl, err := template.New(tmplname).Funcs(funcMap(registry)).Parse(tmplContent)
	if err != nil {
		return fmt.Errorf("parsing template `"+tmplname+"`: %w", err)
	}

	sharedContent, err := tmpldata.Read(SharedTemplateName)
	if err != nil {
		return fmt.Errorf("reading template `"+SharedTemplateName+"`: %w", err)
	}
	gotmpl, err = gotmpl.Parse(sharedContent)
	if err != nil {
		return fmt.Errorf("parsing template `"+SharedTemplateName+"`: %w", err)
	}

	err = gotmpl.Execute(outfile, data)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	if err := gofmt(outpath); err != nil {
		return fmt.Errorf("gofmt: %w", err)
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
