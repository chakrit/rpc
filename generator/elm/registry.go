package elm

type Registry map[string]*ElmType

func (r Registry) Register(t *ElmType) {
	qualifier := t.Module.Name + "." + t.Name
	r[qualifier] = t
}

func (r Registry) Lookup(context *Module, name string) *ElmType {
	qualifier := context.Name + "." + name
	if typ, ok := r[qualifier]; ok {
		return typ
	} else if context.Parent != nil {
		return r.Lookup(context.Parent, name)
	} else {
		return nil
	}
}

// TODO: time, map, etc.
func (r Registry) Resolve(ref *ElmTypeRef) *ElmTypeResolution {
	switch ref.Name {
	case "string", "bool", "int", "long", "float", "double", "time":
		return r.resolveBasic(ref)
	case "list":
		return r.resolveList(ref)
	default:
		return r.resolveUserDefined(ref)
	}
}

func (r Registry) resolveBasic(ref *ElmTypeRef) *ElmTypeResolution {
	switch ref.Name {
	case "string":
		return &ElmTypeResolution{
			Name:   "String",
			Encode: "E.string",
			Decode: "D.string",
		}
	case "bool":
		return &ElmTypeResolution{
			Name:   "Bool",
			Encode: "E.bool",
			Decode: "D.bool",
		}
	case "int", "long":
		return &ElmTypeResolution{
			Name:   "Int",
			Encode: "E.int",
			Decode: "D.int",
		}
	case "float", "double":
		return &ElmTypeResolution{
			Name:   "float",
			Encode: "E.float",
			Decode: "D.float",
		}
	case "time":
		return &ElmTypeResolution{
			Name:   "Time.Posix",
			Encode: `(Time.posixToMillis >> E.int)`,
			Decode: `(D.map Time.millisToPosix D.int)`,
		}
	default:
		return r.resolveUnknown()
	}
}

func (r Registry) resolveList(ref *ElmTypeRef) *ElmTypeResolution {
	var elementType *ElmTypeResolution
	if len(ref.Args) > 0 {
		elementType = r.Resolve(ref.Args[0])
	} else {
		elementType = r.resolveUnknown()
	}

	return r.resolveWithDefault("[]", &ElmTypeResolution{
		Name:   "List (" + elementType.Name + ")",
		Encode: "E.list (" + elementType.Encode + ")",
		Decode: "D.list (" + elementType.Decode + ")",
	})
}

func (r Registry) resolveUserDefined(ref *ElmTypeRef) *ElmTypeResolution {
	elmType := r.Lookup(ref.Module, ref.Name)
	if elmType == nil {
		return r.resolveUnknown() // TODO: Output a warning
	}

	resolved := &ElmTypeResolution{
		Name:   ref.Name,
		Encode: "encode" + ref.Name,
		Decode: "decode" + ref.Name,
	}
	if elmType.Module != ref.Module { // imported type, add module prefix
		resolved.Name = elmType.Module.Name + "." + resolved.Name
		resolved.Encode = elmType.Module.Name + "." + resolved.Encode
		resolved.Decode = elmType.Module.Name + "." + resolved.Decode
	}

	return resolved
}

// tries to give a default value when `null` is received without having to define
// all json-nullable fields as a Maybe
func (r Registry) resolveWithDefault(defaultValue string, elementType *ElmTypeResolution) *ElmTypeResolution {
	dec := "Maybe.withDefault (" + defaultValue + ")"
	dec = "D.map (" + dec + ")"
	dec = dec + " (D.maybe (" + elementType.Decode + "))"

	return &ElmTypeResolution{
		Name:   elementType.Name,
		Encode: elementType.Encode,
		Decode: dec,
	}
}

func (r Registry) resolveUnknown() *ElmTypeResolution {
	return &ElmTypeResolution{
		Name:   "()",
		Encode: `(\_ -> E.null)`,
		Decode: `(D.succeed ())`,
	}
}
