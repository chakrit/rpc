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
	Name    string
	RPCPath string
	OutPath string

	Namespace *spec.Namespace
	Registry  Registry

	Types    []*Type
	Enums    []*Enum
	Tuples   []*Tuple
	RPCFuncs []*RpcFunc
	Imports  []*Module

	Parent   *Module
	Children []*Module
}

func newUtilModule(rootModule *Module, outdir string, ns *spec.Namespace) *Module {
	return &Module{
		Name:      "RpcUtil",
		Namespace: ns,
		Parent:    nil,
		OutPath:   filepath.Join(filepath.Dir(rootModule.OutPath), "RpcUtil.elm"),
	}
}

func newModule(parent *Module, outdir string, ns *spec.Namespace) *Module {
	mod := &Module{
		Namespace: ns,
		Parent:    parent,
	}

	name := ns.Name
	if nameOpt, isSet := ns.Options[ModuleOption]; isSet {
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
		elmType := &Type{
			Name:   typ.Name,
			Module: m,
		}

		for _, p := range typ.Properties.SortedByName() {
			prop := p.(*spec.Property)
			elmType.Fields = append(elmType.Fields, &Field{
				Name: prop.Name,
				Type: m.mapTypeRef(prop.Type),
			})
		}

		m.Types = append(m.Types, elmType)
		m.Registry.RegisterType(elmType)
	}

	for _, e := range m.Namespace.Enums.SortedByName() {
		enum := e.(*spec.Enum)
		elmEnum := &Enum{
			Name:   enum.Name,
			Module: m,
		}

		for _, m := range enum.Members {
			elmEnum.Members = append(elmEnum.Members, &Member{
				Name:  m,
				Value: internal.InflectDash(m),
				Title: internal.InflectTitle(m),
			})
		}

		m.Enums = append(m.Enums, elmEnum)
		m.Registry.RegisterEnum(elmEnum)
	}
}

func (m *Module) resolveRPCFuncs() {
	for _, r := range m.Namespace.RPCs.SortedByName() {
		var (
			rpc    = r.(*spec.RPC)
			inTup  = &Tuple{Name: "InputFor" + rpc.Name}
			outTup = &Tuple{Name: "OutputFor" + rpc.Name}
		)

		for _, ref := range rpc.InputTypes {
			inTup.Args = append(inTup.Args, m.mapTypeRef(ref))
		}
		for _, ref := range rpc.OutputTypes {
			outTup.Args = append(outTup.Args, m.mapTypeRef(ref))
		}

		m.Tuples = append(m.Tuples, inTup, outTup)
		m.RPCFuncs = append(m.RPCFuncs, &RpcFunc{
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
	for _, enum := range m.Enums {
		locals[enum.Name] = struct{}{}
	}

	var check func(ref *TypeRef)
	check = func(ref *TypeRef) {
		for _, arg := range ref.Args {
			check(arg)
		}
		if _, isLocal := locals[ref.Name]; isLocal {
			return // local type, no need to import
		}

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

func (m *Module) mapTypeRef(ref *spec.TypeRef) *TypeRef {
	elmRef := &TypeRef{
		Name:   ref.Name,
		Module: m,
	}

	for _, arg := range ref.Arguments {
		elmRef.Args = append(elmRef.Args, m.mapTypeRef(arg))
	}

	return elmRef
}
