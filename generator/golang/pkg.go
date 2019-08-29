package golang

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/chakrit/rpc/internal"
	"github.com/chakrit/rpc/spec"
)

type pkgByName []*Pkg

func (p pkgByName) Len() int      { return len(p) }
func (p pkgByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p pkgByName) Less(i, j int) bool {
	left, right := p[i], p[j]
	return left.Name < right.Name
}

type PkgContext struct {
	ContextPkg *Pkg
	DataPkg    *Pkg
}

type Pkg struct {
	Name        string
	MangledName string

	BasePath   string
	RPCPath    string
	FilePath   string
	ImportPath string

	Namespace *spec.Namespace
	Registry  TypeRegistry

	Parent   *Pkg
	Children []*Pkg
	Imports  []*Pkg
}

func newRootPkg(ns *spec.Namespace) *Pkg {
	reg := TypeRegistry{}
	pkg := &Pkg{
		Name:      strings.ToLower(ns.Name),
		Namespace: ns,
		Registry:  reg,
	}

	for _, node := range ns.Children.SortedByName() {
		child := newChildPkg(pkg, node.(*spec.Namespace))
		pkg.Children = append(pkg.Children, child)
	}

	sort.Sort(pkgByName(pkg.Children))
	pkg.initialize()
	return pkg
}

func newChildPkg(parent *Pkg, ns *spec.Namespace) *Pkg {
	pkg := &Pkg{
		Name:      strings.ToLower(ns.Name),
		Parent:    parent,
		Namespace: ns,
		Registry:  parent.Registry,
	}
	for _, node := range ns.Children {
		child := newChildPkg(pkg, node.(*spec.Namespace))
		pkg.Children = append(pkg.Children, child)
	}

	sort.Sort(pkgByName(pkg.Children))
	return pkg
}

func (pkg *Pkg) initialize() {
	pkgNameOption, importOption := pkg.lookupOption(PackageOption), pkg.lookupOption(ImportOption)
	if pkgNameOption == "" {
		pkgNameOption = DefaultPkgName
	}
	if importOption == "" {
		importOption = DefaultImportPath
	}
	if pkg.Namespace.Name == "" || pkg.Namespace.Name == "root" {
		pkg.Name = pkgNameOption
	} else {
		pkg.Name = strings.ToLower(pkg.Namespace.Name)
	}

	pkg.resolvePaths("", importOption)
	pkg.generateNames()
	pkg.Registry.RegisterAll(pkg)
	pkg.resolveImports()
}

func (pkg *Pkg) generateNames() {
	if pkg.Parent != nil {
		pkg.MangledName = "rpc_" + internal.InflectSnake(pkg.BasePath)
	} else {
		pkg.MangledName = "rpc_root"
	}

	for _, child := range pkg.Children {
		child.generateNames()
	}
}

func (pkg *Pkg) lookupOption(name string) string {
	if pkg == nil {
		return ""
	} else if value, ok := pkg.Namespace.Options[name]; ok {
		return fmt.Sprint(value)
	} else {
		return pkg.Parent.lookupOption(name)
	}
}

func (pkg *Pkg) resolvePaths(base, importBase string) {
	if pkg.Parent != nil {
		base = path.Join(base, strings.ToLower(pkg.Name))
	}

	pkg.BasePath = base
	pkg.FilePath = path.Join(base, OutName)
	if pkg.Parent != nil {
		pkg.ImportPath = path.Join(importBase, strings.ToLower(pkg.Name))
		pkg.RPCPath = path.Join(pkg.Parent.RPCPath, internal.InflectSnake(pkg.Name))
	} else {
		pkg.ImportPath = importBase
		pkg.RPCPath = path.Join(base, internal.InflectSnake(pkg.Name))
	}

	for _, child := range pkg.Children {
		child.resolvePaths(base, pkg.ImportPath)
	}
}

func (pkg *Pkg) resolveImports() {
	dependencies := map[*Pkg]struct{}{}
	check := func(ref *spec.TypeRef) {
		resolved := pkg.Registry.Resolve(pkg, ref)
		if resolved == nil {
			return
		}

		if resolved.ImportPkg() != nil && resolved.ImportPkg() != pkg {
			dependencies[resolved.ImportPkg()] = struct{}{}
		}
		for _, arg := range resolved.Args() {
			if arg != nil && arg.ImportPkg() != nil && arg.ImportPkg() != pkg {
				dependencies[arg.ImportPkg()] = struct{}{}
			}
		}
	}

	for _, typNode := range pkg.Namespace.Types {
		for _, propNode := range typNode.(*spec.Type).Properties {
			check(propNode.(*spec.Property).Type)
		}
	}
	for _, rpcNode := range pkg.Namespace.RPCs {
		for _, typ := range rpcNode.(*spec.RPC).InputTypes {
			check(typ)
		}
		for _, typ := range rpcNode.(*spec.RPC).OutputTypes {
			check(typ)
		}
	}
	for dependency := range dependencies {
		pkg.Imports = append(pkg.Imports, dependency)
	}

	sort.Sort(pkgByName(pkg.Imports))
	for _, child := range pkg.Children {
		child.resolveImports()
	}
}
