package golang

import (
	"github.com/chakrit/rpc/spec"
)

/*
	"string",
	"bool",
	"int",
	"long",
	"float",
	"double",
	"list",
	"map",
	"time",
*/

type TypeRegistry map[string]*ResolvedType

var (
	resolvedTimeType    *ResolvedType
	resolvedUnknownType *ResolvedType
)

func init() {
	pkg := &Pkg{Name: "time", MangledName: "time", ImportPath: "time"}
	resolvedTimeType = &ResolvedType{"time", pkg, nil, nil, TypeTime | TypeIsImported}
	resolvedUnknownType = &ResolvedType{"interface{}", nil, nil, nil, TypeSimple | TypeIsBuiltin}
}

func (r TypeRegistry) RegisterAll(pkg *Pkg) {
	for _, typNode := range pkg.Namespace.Types {
		typ := typNode.(*spec.Type)
		slug := r.slug(pkg, typ.Name)
		r[slug] = &ResolvedType{
			Name:  typ.Name,
			Pkg:   pkg,
			Type:  typ,
			Flags: TypeUserDefined,
		}
	}

	for _, child := range pkg.Children {
		r.RegisterAll(child)
	}
}

func (r TypeRegistry) Resolve(pkg *Pkg, ref *spec.TypeRef) *ResolvedType {
	switch ref.Name {
	case "string", "bool", "int", "long", "float", "double", "time":
		return r.resolveSimpleType(pkg, ref)
	case "list":
		return r.resolveListType(pkg, ref)
	case "map":
		return r.resolveMapType(pkg, ref)
	default:
		return r.resolveUserDefinedType(pkg, ref)
	}
}

func (r TypeRegistry) resolveSimpleType(pkg *Pkg, ref *spec.TypeRef) *ResolvedType {
	switch ref.Name {
	case "string", "bool", "int":
		return &ResolvedType{ref.Name, nil, nil, nil, TypeSimple | TypeIsBuiltin}
	case "long":
		return &ResolvedType{"int64", nil, nil, nil, TypeSimple | TypeIsBuiltin}
	case "float":
		return &ResolvedType{"float32", nil, nil, nil, TypeSimple | TypeIsBuiltin}
	case "double":
		return &ResolvedType{"float64", nil, nil, nil, TypeSimple | TypeIsBuiltin}
	case "time":
		return resolvedTimeType
	default:
		return resolvedUnknownType
	}
}

func (r TypeRegistry) resolveListType(pkg *Pkg, ref *spec.TypeRef) *ResolvedType {
	switch len(ref.Arguments) {
	case 0:
		// TODO: Emit a warning, list with no type arg
		args := []*ResolvedType{resolvedUnknownType}
		return &ResolvedType{"[]", nil, nil, args, TypeList | TypeIsBuiltin}
	case 1:
		args := []*ResolvedType{r.Resolve(pkg, ref.Arguments[0])}
		return &ResolvedType{"[]", nil, nil, args, TypeList | TypeIsBuiltin}
	default:
		// TODO: Emit a warning, list with too many type args
		args := []*ResolvedType{resolvedUnknownType}
		return &ResolvedType{"[]", nil, nil, args, TypeList | TypeIsBuiltin}
	}
}

func (r TypeRegistry) resolveMapType(pkg *Pkg, ref *spec.TypeRef) *ResolvedType {
	switch len(ref.Arguments) {
	case 0:
		// TODO: Emit a warning, map with no key and value type args
		args := []*ResolvedType{
			resolvedUnknownType,
			resolvedUnknownType,
		}
		return &ResolvedType{"map", nil, nil, args, TypeMap | TypeIsBuiltin}
	case 1:
		args := []*ResolvedType{
			r.Resolve(pkg, ref.Arguments[0]),
			resolvedUnknownType,
		}
		return &ResolvedType{"map", nil, nil, args, TypeMap | TypeIsBuiltin}
	case 2:
		args := []*ResolvedType{
			r.Resolve(pkg, ref.Arguments[0]),
			r.Resolve(pkg, ref.Arguments[1]),
		}
		return &ResolvedType{"map", nil, nil, args, TypeMap | TypeIsBuiltin}
	default:
		// TODO: Emit a warning, map with too many arguments types
		args := []*ResolvedType{
			r.Resolve(pkg, ref.Arguments[0]),
			r.Resolve(pkg, ref.Arguments[1]),
		}
		return &ResolvedType{"map", nil, nil, args, TypeMap | TypeIsBuiltin}
	}
}

func (r TypeRegistry) resolveUserDefinedType(pkg *Pkg, ref *spec.TypeRef) *ResolvedType {
	for findPkg := pkg; findPkg != nil; findPkg = findPkg.Parent {
		slug := r.slug(findPkg, ref.Name)
		if existing, ok := r[slug]; ok {
			if findPkg == pkg {
				return existing.WithFlags(TypeIsLocal)
			} else {
				return existing.WithFlags(TypeIsImported)
			}
		}
	}

	// TODO: Emit a warning, can't find user-defind type
	return resolvedUnknownType
}

func (r TypeRegistry) slug(pkg *Pkg, name string) string {
	return pkg.BasePath + "." + name
}
