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

func (r Registry) Resolve(ref *ElmTypeRef) *ElmTypeResolution {
	switch ref.Name {
	case "string", "bool", "int", "long", "float", "double", "time", "data":
		return r.resolveBasic(ref)
	case "list":
		return r.resolveList(ref)
	case "map":
		return r.resolveMap(ref)
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
			Name:   "Posix",
			Encode: `(Time.posixToMillis >> toFloat >> (\f -> f/1000.0) >> E.float)`,
			Decode: `(D.map ((\f -> f * 1000.0) >> round >> Time.millisToPosix) D.float)`,
		}
	case "data":
		return &ElmTypeResolution{
			Name:   "Bytes",
			Encode: `(RpcUtil.b64StringFromBytes >> Maybe.withDefault "" >> E.string)`,
			Decode: `(D.map (Maybe.withDefault "" >> RpcUtil.b64StringToBytes >> Maybe.withDefault (Bytes.Encode.encode (Bytes.Encode.string ""))) (D.maybe D.string))`,
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

func (r Registry) resolveMap(ref *ElmTypeRef) *ElmTypeResolution {
	var keyType *ElmTypeResolution
	var valueType *ElmTypeResolution
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
	return r.resolveWithDefault("Dict.empty", &ElmTypeResolution{
		Name:   "Dict (" + keyType.Name + ") (" + valueType.Name + ")",
		Encode: "E.dict (identity) (" + valueType.Encode + ")",
		Decode: "D.dict (" + valueType.Decode + ")",
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
//
// TODO: Move this logic into the template and force all types to have defaults?
func (r Registry) resolveWithDefault(defaultValue string, elementType *ElmTypeResolution) *ElmTypeResolution {
	// i.e. D.map (Maybe.withDefault ([])) (D.maybe (D.list decodeSomething))
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
