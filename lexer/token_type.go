package lexer

import "encoding/json"

const (
	T_Unknown = TokenType(1 << iota)
	T_Space
	T_EndOfLine
	T_EndOfFile
	T_BlockStart
	T_BlockEnd
	T_ArgListStart
	T_ArgListEnd
	T_ArgListSep
	T_TypeArgListStart
	T_TypeArgListEnd
	T_Comment
	T_Identifier
	T_Keyword
	T_StringValue
	T_NumberValue
)

var typeNames = map[TokenType]string{
	T_Unknown:          "unknown",
	T_Space:            "whitespace",
	T_EndOfLine:        "end-of-line",
	T_EndOfFile:        "end-of-file",
	T_BlockStart:       "block-start",
	T_BlockEnd:         "block-end",
	T_ArgListStart:     "arg-list-start",
	T_ArgListEnd:       "arg-list-end",
	T_ArgListSep:       "arg-list-sep",
	T_TypeArgListStart: "type-arg-list-start",
	T_TypeArgListEnd:   "type-arg-list-end",
	T_Comment:          "comment",
	T_Identifier:       "identifier",
	T_Keyword:          "keyword",
	T_StringValue:      "value-string",
	T_NumberValue:      "value-number",
}

var braceMappings = map[rune]TokenType{
	'{': T_BlockStart,
	'}': T_BlockEnd,
	'(': T_ArgListStart,
	')': T_ArgListEnd,
	'<': T_TypeArgListStart,
	'>': T_TypeArgListEnd,
}

type TokenType int

func (t TokenType) Match(another TokenType) bool {
	return (int(t) & int(another)) > 0
}

func (t TokenType) String() string {
	if name, ok := typeNames[t]; ok {
		return name
	} else {
		return typeNames[T_Unknown]
	}
}

func (t TokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
