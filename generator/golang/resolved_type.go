package golang

var (
	unknownType ResolvedType = rtSimple{"interface{}"}
	unitType    ResolvedType = rtSimple{"struct{}"}
	stringType  ResolvedType = rtSimple{"string"}
	boolType    ResolvedType = rtSimple{"bool"}
	intType     ResolvedType = rtSimple{"int"}
	longType    ResolvedType = rtSimple{"int64"}
	floatType   ResolvedType = rtSimple{"float32"}
	doubleType  ResolvedType = rtSimple{"float64"}
	dataType    ResolvedType = rtSimple{"[]byte"}

	timeType ResolvedType = rtTime{}
)

type (
	ResolvedType interface {
		Name() string
		Args() []ResolvedType
		ImportPkg() *Pkg

		AsReference(current *Pkg) string
	}

	CustomMarshaler interface {
		AsMarshalTarget(current *Pkg) string
		AsMarshaler(current *Pkg) string
		AsUnmarshaler(current *Pkg) string
	}

	rtSimple struct{ name string }
	rtTime   struct{}
	rtList   struct{ arg ResolvedType }
	rtMap    struct {
		keyArg   ResolvedType
		valueArg ResolvedType
	}

	rtEnum struct {
		name      string
		importPkg *Pkg
	}
	rtUserDefined struct {
		name      string
		importPkg *Pkg
	}
)

func (t rtSimple) Name() string                { return t.name }
func (t rtSimple) Args() []ResolvedType        { return nil }
func (t rtSimple) ImportPkg() *Pkg             { return nil }
func (t rtSimple) AsReference(cur *Pkg) string { return t.name }

func (t rtTime) Name() string         { return "time" }
func (t rtTime) Args() []ResolvedType { return nil }
func (t rtTime) ImportPkg() *Pkg {
	return &Pkg{
		Name:        "time",
		MangledName: "time",
		ImportPath:  "time",
	}
}
func (t rtTime) AsReference(cur *Pkg) string { return "time.Time" }
func (t rtTime) AsMarshalTarget(cur *Pkg) string {
	return "float64"
}
func (t rtTime) AsMarshaler(current *Pkg) string {
	return "(func (t time.Time) float64 {" +
		"sec, nsec := t.Unix(), t.Nanosecond();" +
		"return float64(sec) + (float64(nsec)/float64(time.Second));" +
		"})"
}
func (t rtTime) AsUnmarshaler(current *Pkg) string {
	return "(func (t float64) time.Time {" +
		"fsec, fnsec := math.Modf(t);" +
		"sec, nsec := int64(fsec), int64(math.Round(fnsec*float64(time.Second)));" +
		"return time.Unix(sec, nsec);" +
		"})"
}

func (t rtList) Name() string         { return "[]" }
func (t rtList) Args() []ResolvedType { return []ResolvedType{t.arg} }
func (t rtList) ImportPkg() *Pkg      { return nil }
func (t rtList) AsReference(cur *Pkg) string {
	return "[]" + t.arg.AsReference(cur)
}

func (t rtMap) Name() string         { return "map" }
func (t rtMap) Args() []ResolvedType { return []ResolvedType{t.keyArg, t.valueArg} }
func (t rtMap) ImportPkg() *Pkg      { return nil }
func (t rtMap) AsReference(cur *Pkg) string {
	return "map[" + t.keyArg.AsReference(cur) +
		"]" + t.valueArg.AsReference(cur)
}

func (t rtEnum) Name() string         { return t.name }
func (t rtEnum) Args() []ResolvedType { return nil }
func (t rtEnum) ImportPkg() *Pkg      { return t.importPkg }
func (t rtEnum) AsReference(cur *Pkg) string {
	if cur == t.importPkg {
		return t.name
	} else {
		return t.importPkg.MangledName + "." + t.name
	}
}
func (t rtEnum) AsMarshalTarget(cur *Pkg) string {
	return "string"
}
func (t rtEnum) AsMarshaler(cur *Pkg) string {
	ref := t.AsReference(cur)
	return "(func(v " + ref + ") string { return string(v) })"
}
func (t rtEnum) AsUnmarshaler(cur *Pkg) string {
	ref := t.AsReference(cur)
	return "(func(v string) " + ref + " { return " + ref + "(v) })"
}

func (t rtUserDefined) Name() string         { return t.name }
func (t rtUserDefined) Args() []ResolvedType { return nil }
func (t rtUserDefined) ImportPkg() *Pkg      { return t.importPkg }
func (t rtUserDefined) AsReference(cur *Pkg) string {
	if cur == t.importPkg {
		return "*" + t.name
	} else {
		return "*" + t.importPkg.MangledName + "." + t.name
	}
}
