package lexer

var Keywords = map[string]struct{}{
	// meta
	"option":  {},
	"include": {},
	"root":    {}, // unused, reserved as root ns identifier

	// definitions
	"namespace": {},
	"type":      {},
	"rpc":       {},

	// built-in types
	"string": {},
	"bool":   {},
	"int":    {},
	"long":   {},
	"float":  {},
	"double": {},
	"list":   {},
	"map":    {},
	"time":   {},
	"data":   {},
}

func IsKeyword(word string) bool {
	_, ok := Keywords[word]
	return ok
}
