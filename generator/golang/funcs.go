package golang

import (
	"text/template"

	"github.com/chakrit/rpc/internal"
	"github.com/chakrit/rpc/spec"
)

func funcMap(pkg *Pkg) template.FuncMap {
	reg, f := pkg.Registry, template.FuncMap{}

	f["pascal"] = internal.InflectPascal
	f["snake"] = internal.InflectSnake

	f["asReference"] = AsReference

	f["resolve"] = reg.Resolve
	f["resolveAbs"] = func(pkg *Pkg, ref *spec.TypeRef) *ResolvedType {
		if resolved := reg.Resolve(pkg, ref); resolved == nil {
			return nil
		} else {
			return resolved.
				WithoutFlags(OriginPkgLocal).
				WithFlags(OriginImported)
		}
	}

	return f
}

func AsReference(rt *ResolvedType) string {
	if rt == nil {
		// TODO: Warn about bad resolution
		return "interface{}"
	}

	var (
		arg0 *ResolvedType
		arg1 *ResolvedType
	)

	if len(rt.Args) > 1 {
		arg1 = rt.Args[1]
	}
	if len(rt.Args) > 0 {
		arg0 = rt.Args[0]
	}

	switch {
	case rt.IsSimple():
		return rt.Name
	case rt.IsTime():
		return "time.Time"
	case rt.IsMap():
		return "map[" + AsReference(arg0) + "]" + AsReference(arg1)
	case rt.IsList():
		return "[]" + AsReference(arg0)
	case rt.IsPkgLocal():
		return "*" + rt.Name
	case rt.IsImported():
		return "*" + rt.Pkg.MangledName + "." + rt.Name
	default:
		// TODO: Warn about bad resolution
		return "interface{}"
	}

}
