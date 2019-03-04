package golang

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/chakrit/rpc/internal"
	"github.com/chakrit/rpc/spec"
)

type pkgByNameAndNumber []*Pkg

func (p pkgByNameAndNumber) Len() int      { return len(p) }
func (p pkgByNameAndNumber) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p pkgByNameAndNumber) Less(i, j int) bool {
	left, right := p[i], p[j]
	return left.Name < right.Name ||
		left.Number < right.Number
}

type Pkg struct {
	Number       int // increasing number, for various uniqueness requirements
	Name         string
	ExportedName string
	MangledName  string

	RPCPath    string
	FilePath   string
	ImportPath string

	Namespace *spec.Namespace
	Registry  TypeRegistry
	Funcs     template.FuncMap

	Parent       *Pkg
	Children     []*Pkg
	Dependencies []*Pkg
}

func newRootPkg(ns *spec.Namespace) *Pkg {
	reg := TypeRegistry{}
	pkg := &Pkg{
		Name:      ns.Name,
		Namespace: ns,
		Registry:  reg,
	}

	pkg.Funcs = funcMap(pkg)
	for _, node := range ns.Children.SortedByName() {
		child := newChildPkg(pkg, node.(*spec.Namespace))
		pkg.Children = append(pkg.Children, child)
	}

	sort.Sort(pkgByNameAndNumber(pkg.Children))
	pkg.initialize()
	return pkg
}

func newChildPkg(parent *Pkg, ns *spec.Namespace) *Pkg {
	pkg := &Pkg{
		Name:      ns.Name,
		Parent:    parent,
		Namespace: ns,
		Registry:  parent.Registry,
		Funcs:     parent.Funcs,
	}
	for _, node := range ns.Children {
		child := newChildPkg(pkg, node.(*spec.Namespace))
		pkg.Children = append(pkg.Children, child)
	}

	sort.Sort(pkgByNameAndNumber(pkg.Children))
	return pkg
}

func (pkg *Pkg) initialize() {
	pkg.assignNamesAndNumbers(nil)
	pkgNameOption := pkg.lookupOption(PackageOption)
	if pkg.Name == "" || pkg.Name == "root" {
		pkg.Name = pkgNameOption
	}

	pkg.resolvePaths("", pkg.lookupOption(ImportOption))
	pkg.Registry.RegisterAll(pkg)
	pkg.resolveImports()
	pkg.generateAltNames()
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

func (pkg *Pkg) assignNamesAndNumbers(count *int) {
	if count == nil {
		n := 1
		count = &n
	} else {
		*count += 1
	}

	pkg.Number = *count
	pkg.Name = strings.ToLower(pkg.Namespace.Name)
	for _, child := range pkg.Children {
		child.assignNamesAndNumbers(count)
	}
}

func (pkg *Pkg) resolvePaths(base, importBase string) {
	if pkg.Parent != nil {
		base = path.Join(base, strings.ToLower(pkg.Name))
	}

	pkg.RPCPath = path.Join(base, internal.InflectSnake(pkg.Name))
	pkg.FilePath = path.Join(base, OutName)
	if pkg.Parent != nil {
		pkg.ImportPath = path.Join(importBase, strings.ToLower(pkg.Name))
	} else {
		pkg.ImportPath = importBase
	}

	for _, child := range pkg.Children {
		child.resolvePaths(base, pkg.ImportPath)
	}
}

func (pkg *Pkg) resolveImports() {
	dependencies := map[*Pkg]bool{}
	check := func(ref *spec.TypeRef) {
		resolved := pkg.Registry.Resolve(pkg, ref)

		// TODO: Warn if resolved == nil
		if resolved != nil &&
			resolved.Pkg != nil &&
			resolved.Pkg != pkg {
			dependencies[resolved.Pkg] = true
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
		pkg.Dependencies = append(pkg.Dependencies, dependency)
	}

	sort.Sort(pkgByNameAndNumber(pkg.Dependencies))
	for _, child := range pkg.Children {
		child.resolveImports()
	}
}

func (pkg *Pkg) generateAltNames() {
	pkg.ExportedName = internal.InflectPascal(pkg.Name)
	pkg.MangledName = internal.MangleID(
		pkg.ImportPath,
		strconv.Itoa(pkg.Number))

	for _, child := range pkg.Children {
		child.generateAltNames()
	}
}
