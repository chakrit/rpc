package golang

import (
	"github.com/chakrit/rpc/spec"
)

type TypeRegistry map[string]ResolvedType

func (r TypeRegistry) RegisterAll(pkg *Pkg) {
	for _, typNode := range pkg.Namespace.Types {
		typ := typNode.(*spec.Type)
		slug := r.slug(pkg, typ.Name)
		r[slug] = rtUserDefined{typ.Name, pkg}
	}
	for _, enumNode := range pkg.Namespace.Enums {
		enum := enumNode.(*spec.Enum)
		slug := r.slug(pkg, enum.Name)
		r[slug] = rtEnum{enum.Name, pkg}
	}

	for _, child := range pkg.Children {
		r.RegisterAll(child)
	}
}

func (r TypeRegistry) Resolve(pkg *Pkg, ref *spec.TypeRef) ResolvedType {
	switch ref.Name {
	case "unit", "string", "bool", "int", "long", "float", "double", "time", "data":
		return r.resolveSimpleType(pkg, ref)
	case "list":
		return r.resolveListType(pkg, ref)
	case "map":
		return r.resolveMapType(pkg, ref)
	default:
		return r.resolveCustomType(pkg, ref)
	}
}

func (r TypeRegistry) resolveSimpleType(pkg *Pkg, ref *spec.TypeRef) ResolvedType {
	switch ref.Name {
	case "unit":
		return unitType
	case "string":
		return stringType
	case "bool":
		return boolType
	case "int":
		return intType
	case "long":
		return longType
	case "float":
		return floatType
	case "double":
		return doubleType
	case "data":
		return dataType
	case "time":
		return timeType
	default:
		return unknownType
	}
}

func (r TypeRegistry) resolveListType(pkg *Pkg, ref *spec.TypeRef) ResolvedType {
	switch len(ref.Arguments) {
	case 0:
		// TODO: Emit a warning, list with no type arg
		return rtList{unknownType}
	case 1:
		return rtList{r.Resolve(pkg, ref.Arguments[0])}
	default:
		// TODO: Emit a warning, list with too many type args
		return rtList{r.Resolve(pkg, ref.Arguments[0])}
	}
}

func (r TypeRegistry) resolveMapType(pkg *Pkg, ref *spec.TypeRef) ResolvedType {
	switch len(ref.Arguments) {
	case 0:
		// TODO: Emit a warning, map with no key and value type args
		return rtMap{unknownType, unknownType}
	case 1:
		return rtMap{r.Resolve(pkg, ref.Arguments[0]), unknownType}
	case 2:
		return rtMap{
			r.Resolve(pkg, ref.Arguments[0]),
			r.Resolve(pkg, ref.Arguments[1]),
		}
	default:
		// TODO: Emit a warning, map with too many arguments types
		return rtMap{
			r.Resolve(pkg, ref.Arguments[0]),
			r.Resolve(pkg, ref.Arguments[1]),
		}
	}
}

func (r TypeRegistry) resolveCustomType(pkg *Pkg, ref *spec.TypeRef) ResolvedType {
	for findPkg := pkg; findPkg != nil; findPkg = findPkg.Parent {
		slug := r.slug(findPkg, ref.Name)
		switch r[slug].(type) {
		case nil:
			continue
		case rtUserDefined:
			return rtUserDefined{ref.Name, findPkg}
		case rtEnum:
			return rtEnum{ref.Name, findPkg}
		}
	}

	// TODO: Emit a warning, can't find user-defined type
	return unknownType
}

func (r TypeRegistry) slug(pkg *Pkg, name string) string {
	return pkg.BasePath + "." + name
}
