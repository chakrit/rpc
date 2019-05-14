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
