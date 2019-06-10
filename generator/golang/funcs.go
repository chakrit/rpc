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

	f["asReference"] = asReference
	f["asMarshalTarget"] = asMarshalTarget
	f["marshaler"] = marshaler
	f["unmarshaler"] = unmarshaler

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

func asReference(rt *ResolvedType) string {
	if rt == nil {
		// TODO: Warn about bad resolution
		return "interface{}"
	}

	arg0, arg1 := rt.Arg(0), rt.Arg(1)

	switch {
	case rt.IsSimple():
		return rt.Name
	case rt.IsTime():
		return "time.Time"
	case rt.IsMap():
		return "map[" + asReference(arg0) + "]" + asReference(arg1)
	case rt.IsList():
		return "[]" + asReference(arg0)
	case rt.IsPkgLocal():
		return "*" + rt.Name
	case rt.IsImported():
		return "*" + rt.Pkg.MangledName + "." + rt.Name
	default:
		// TODO: Warn about bad resolution
		return "interface{}"
	}

}

func asMarshalTarget(rt *ResolvedType) string {
	if rt == nil {
		// TODO: Warn about bad resolution
		return "interface{}"
	}

	switch {
	case rt.IsTime():
		// Posix are more universal than iso8601 -> this allows more language to parse and
		// encode it without having to add a third party dependency. We can still retain
		// some precision (not perfect, but good enough) from go's nanoseconds by going
		// into the decimals.
		return "float64"
	default:
		return asReference(rt)
	}
}

// TODO: Move into ResolvedType
func marshaler(rt *ResolvedType) string {
	if rt == nil {
		return ""
	}

	switch {
	case rt.IsTime():
		return "(func (t time.Time) float64 {" +
			"sec, nsec := t.Unix(), t.Nanosecond();" +
			"return float64(sec) + (float64(nsec)/float64(time.Second));" +
			"})"
	case rt.IsBytes():
		return "(func (b []byte) string { return string(b) })"
	default:
		return ""
	}
}

func unmarshaler(rt *ResolvedType) string {
	if rt == nil {
		return ""
	}

	switch {
	case rt.IsTime():
		return "(func (t float64) time.Time {" +
			"fsec, fnsec := math.Modf(t);" +
			"sec, nsec := int64(fsec), int64(math.Round(fnsec*float64(time.Second)));" +
			"return time.Unix(sec, nsec);" +
			"})"
	case rt.IsBytes():
		return "(func (s string) []byte { return []byte(s) })"
	default:
		return ""
	}
}
