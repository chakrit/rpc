package golang

import (
	"text/template"

	"github.com/chakrit/rpc/internal"
)

func funcMap(reg TypeRegistry) template.FuncMap {
	f := template.FuncMap{}
	f["pascal"] = internal.InflectPascal
	f["snake"] = internal.InflectSnake

	f["context"] = tmplContext

	f["resolve"] = reg.Resolve
	f["asReference"] = asReference
	f["asMarshalTarget"] = asMarshalTarget
	f["asMarshaler"] = asMarshaler
	f["asUnmarshaler"] = asUnmarshaler
	return f
}

func tmplContext(ctxPkg *Pkg, dataPkg *Pkg) *PkgContext {
	return &PkgContext{
		ContextPkg: ctxPkg,
		DataPkg:    dataPkg,
	}
}

func asReference(pkg *Pkg, rt ResolvedType) string {
	if rt == nil {
		// TODO: Warn about bad resolution
		rt = unknownType
	}

	return rt.AsReference(pkg)
}

func asMarshalTarget(pkg *Pkg, rt ResolvedType) string {
	if rt == nil {
		// TODO: Warn about bad resolution
		rt = unknownType
	}

	if m, ok := rt.(CustomMarshaler); ok {
		return m.AsMarshalTarget(pkg)
	} else {
		return asReference(pkg, rt)
	}
}

func asMarshaler(pkg *Pkg, rt ResolvedType) string {
	if rt == nil {
		rt = unknownType
	}

	if m, ok := rt.(CustomMarshaler); ok {
		return m.AsMarshaler(pkg)
	} else {
		return ""
	}
}

func asUnmarshaler(pkg *Pkg, rt ResolvedType) string {
	if rt == nil {
		rt = unknownType
	}

	if m, ok := rt.(CustomMarshaler); ok {
		return m.AsUnmarshaler(pkg)
	} else {
		return ""
	}
}
