package elm

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chakrit/rpc/internal"
	"github.com/chakrit/rpc/spec"
)

type Module struct {
	Number  int // increasing number, for uniqueness
	Name    string
	RPCPath string
	OutPath string

	Namespace *spec.Namespace
	Registry  Registry

	Types    []*ElmType
	Tuples   []*ElmTuple
	RPCFuncs []*ElmRPCFunc
	Imports  []*Module

	Parent   *Module
	Children []*Module
}

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

var counter = 0

func newSharedModule(rootModule *Module, outdir string, ns *spec.Namespace) *Module {
	return &Module{
		Name:      "RpcUtil",
		Namespace: ns,
		Parent:    nil,
		OutPath:   filepath.Join(filepath.Dir(rootModule.OutPath), "RpcUtil.elm"),
	}
}

func newModule(parent *Module, outdir string, ns *spec.Namespace) *Module {
	counter += 1
	mod := &Module{
		Number:    counter,
		Namespace: ns,
		Parent:    parent,
	}

	name := ns.Name
	if nameOpt, isSet := ns.Options[ElmModuleOption]; isSet {
		name = fmt.Sprint(nameOpt)
	} else if name == "root" {
		name = "Rpc"
	}

	pascalName := internal.InflectPascal(name)
	lowerName := strings.ToLower(name)
	if parent == nil {
		mod.Name = pascalName
		mod.Registry = Registry{}
		mod.OutPath = filepath.Join(outdir, mod.Name+".elm")
		mod.RPCPath = lowerName
	} else { // parent != nil
		mod.Name = parent.Name + "." + pascalName
		mod.Registry = parent.Registry
		mod.OutPath = filepath.Join(outdir, pascalName+".elm")
		mod.RPCPath = path.Join(parent.RPCPath, lowerName)
	}

	mod.resolveTypes()
	mod.resolveRPCFuncs()
	for _, node := range ns.Children.SortedByName() {
		child := node.(*spec.Namespace)
		childDir := filepath.Join(outdir, pascalName)
		mod.Children = append(mod.Children, newModule(mod, childDir, child))
	}

	mod.resolveImports() // after we have RPCs and Types refs
	return mod
}

func (m *Module) resolveTypes() {
	for _, t := range m.Namespace.Types.SortedByName() {
		typ := t.(*spec.Type)
		elmType := &ElmType{
			Name:   typ.Name,
			Module: m,
		}

		for _, p := range typ.Properties.SortedByName() {
			prop := p.(*spec.Property)
			elmType.Fields = append(elmType.Fields, &ElmField{
				Name: prop.Name,
				Type: m.mapTypeRef(prop.Type),
			})
		}

		m.Types = append(m.Types, elmType)
		m.Registry.Register(elmType)
	}
}

func (m *Module) resolveRPCFuncs() {
	for _, r := range m.Namespace.RPCs.SortedByName() {
		var (
			rpc    = r.(*spec.RPC)
			inTup  = &ElmTuple{Name: "InputFor" + rpc.Name}
			outTup = &ElmTuple{Name: "OutputFor" + rpc.Name}
		)

		for _, ref := range rpc.InputTypes {
			inTup.Args = append(inTup.Args, m.mapTypeRef(ref))
		}
		for _, ref := range rpc.OutputTypes {
			outTup.Args = append(outTup.Args, m.mapTypeRef(ref))
		}

		m.Tuples = append(m.Tuples, inTup, outTup)
		m.RPCFuncs = append(m.RPCFuncs, &ElmRPCFunc{
			Name:    rpc.Name,
			RPCPath: path.Join(m.RPCPath, rpc.Name),
			InArgs:  inTup.Args,
			OutArgs: outTup.Args,
		})
	}
}

func (m *Module) resolveImports() {
	locals := map[string]struct{}{}
	imports := map[string]struct{}{}
	for _, typ := range m.Types {
		locals[typ.Name] = struct{}{}
	}

	check := func(ref *ElmTypeRef) {
		if _, isLocal := locals[ref.Name]; !isLocal {
			typ := m.Registry.Lookup(m, ref.Name)
			if typ == nil {
				// TODO: Emit a warning
				return
			}

			if _, imported := imports[typ.Module.Name]; !imported {
				m.Imports = append(m.Imports, typ.Module)
				imports[typ.Module.Name] = struct{}{}
			}
		}
	}

	for _, typ := range m.Types {
		for _, field := range typ.Fields {
			check(field.Type)
		}
	}
	for _, tup := range m.Tuples {
		for _, arg := range tup.Args {
			check(arg)
		}
	}

	sort.Slice(m.Imports, func(i, j int) bool {
		mi, mj := m.Imports[i], m.Imports[j]
		return mi.Name < mj.Name
	})
}

func (m *Module) mapTypeRef(ref *spec.TypeRef) *ElmTypeRef {
	elmRef := &ElmTypeRef{
		Name:   ref.Name,
		Module: m,
	}

	for _, arg := range ref.Arguments {
		elmRef.Args = append(elmRef.Args, m.mapTypeRef(arg))
	}

	return elmRef
}
