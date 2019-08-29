package elm

type Registry map[string]*Type

func (r Registry) Register(t *Type) {
	qualifier := t.Module.Name + "." + t.Name
	r[qualifier] = t
}

func (r Registry) Lookup(context *Module, name string) *Type {
	qualifier := context.Name + "." + name
	if typ, ok := r[qualifier]; ok {
		return typ
	} else if context.Parent != nil {
		return r.Lookup(context.Parent, name)
	} else {
		return nil
	}
}

func (r Registry) Resolve(ref *TypeRef) *TypeResolution {
	switch ref.Name {
	case "unit", "string", "bool", "int", "long", "float", "double", "time", "data":
		return r.resolveBasic(ref)
	case "list":
		return r.resolveList(ref)
	case "map":
		return r.resolveMap(ref)
	default:
		return r.resolveUserDefined(ref)
	}
}

func (r Registry) resolveBasic(ref *TypeRef) *TypeResolution {
	switch ref.Name {
	case "unit":
		return &TypeResolution{
			Name:   "()",
			Encode: `(\_ -> E.object [])`,
			Decode: `D.map (\_ -> ()) D.value`,
		}
	case "string":
		return &TypeResolution{
			Name:    "String",
			Encode:  "E.string",
			Decode:  "D.string",
			Default: `""`,
		}
	case "bool":
		return &TypeResolution{
			Name:    "Bool",
			Encode:  "E.bool",
			Decode:  "D.bool",
			Default: "False",
		}
	case "int", "long":
		return &TypeResolution{
			Name:    "Int",
			Encode:  "E.int",
			Decode:  "D.int",
			Default: "0",
		}
	case "float", "double":
		return &TypeResolution{
			Name:    "Float",
			Encode:  "E.float",
			Decode:  "D.float",
			Default: "0.0",
		}
	case "time":
		return &TypeResolution{
			Name:    "Posix",
			Encode:  `(Time.posixToMillis >> toFloat >> (\f -> f/1000.0) >> E.float)`,
			Decode:  `(D.map ((\f -> f * 1000.0) >> round >> Time.millisToPosix) D.float)`,
			Default: `Time.millisToPosix 0`,
		}
	case "data":
		// data url is used on most things on Elm-side and we have easy conversion with
		// File.toUrl so we assume simple string handling on Elm side for binary data
		return &TypeResolution{
			Name:    "String",
			Encode:  `E.string`,
			Decode:  `D.string`,
			Default: `""`,
		}
	default:
		return r.resolveUnknown()
	}
}

func (r Registry) resolveList(ref *TypeRef) *TypeResolution {
	var elementType *TypeResolution
	if len(ref.Args) > 0 {
		elementType = r.Resolve(ref.Args[0])
	} else {
		elementType = r.resolveUnknown()
	}

	return &TypeResolution{
		Name:    "List (" + elementType.Name + ")",
		Encode:  "E.list (" + elementType.Encode + ")",
		Decode:  "D.list (" + elementType.Decode + ")",
		Default: "[]",
	}
}

func (r Registry) resolveMap(ref *TypeRef) *TypeResolution {
	var keyType *TypeResolution
	var valueType *TypeResolution
	switch len(ref.Args) {
	case 0:
		keyType, valueType = r.resolveUnknown(), r.resolveUnknown()
	case 1:
		keyType, valueType = r.Resolve(ref.Args[0]), r.resolveUnknown()
	default:
		keyType, valueType = r.Resolve(ref.Args[0]), r.Resolve(ref.Args[1])
	}

	// TODO: Warn for non-string key types not supported. Currently this code
	//   will not compile correctly if the key type is not a string
	return &TypeResolution{
		Name:    "Dict (" + keyType.Name + ") (" + valueType.Name + ")",
		Encode:  "E.dict (identity) (" + valueType.Encode + ")",
		Decode:  "D.dict (" + valueType.Decode + ")",
		Default: "Dict.empty",
	}
}

func (r Registry) resolveUserDefined(ref *TypeRef) *TypeResolution {
	elmType := r.Lookup(ref.Module, ref.Name)
	if elmType == nil {
		return r.resolveUnknown() // TODO: Output a warning
	}

	resolved := &TypeResolution{
		Name:    ref.Name,
		Encode:  "encode" + ref.Name,
		Decode:  "decode" + ref.Name,
		Default: "default" + ref.Name,
	}
	if elmType.Module != ref.Module { // imported type, add module prefix
		resolved.Name = elmType.Module.Name + "." + resolved.Name
		resolved.Encode = elmType.Module.Name + "." + resolved.Encode
		resolved.Decode = elmType.Module.Name + "." + resolved.Decode
		resolved.Default = elmType.Module.Name + "." + resolved.Default
	}

	return resolved
}

func (r Registry) resolveUnknown() *TypeResolution {
	return &TypeResolution{
		Name:   "()",
		Encode: `(\_ -> E.null)`,
		Decode: `(D.succeed ())`,
	}
}
