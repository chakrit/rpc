package lexer

import (
	"fmt"

	"github.com/chakrit/rpc/internal"
)

const (
	T_Space = 1 << iota
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

var typeNames = map[int]string{
	T_Space:            "Whitespace",
	T_EndOfLine:        "End of Line",
	T_EndOfFile:        "End of File",
	T_BlockStart:       "Block Start",
	T_BlockEnd:         "Block End",
	T_ArgListStart:     "Args Start",
	T_ArgListEnd:       "Args End",
	T_ArgListSep:       "Args Sep",
	T_TypeArgListStart: "Types Start",
	T_TypeArgListEnd:   "Types End",
	T_Comment:          "Comment",
	T_Identifier:       "Identifier",
	T_Keyword:          "Keyword",
	T_StringValue:      "Value Str",
	T_NumberValue:      "Value Num",
}

var braceMappings = map[rune]int{
	'{': T_BlockStart,
	'}': T_BlockEnd,
	'(': T_ArgListStart,
	')': T_ArgListEnd,
	'<': T_TypeArgListStart,
	'>': T_TypeArgListEnd,
}

type Token struct {
	Type  int          `json:"type"`
	Value string       `json:"value"`
	Pos   internal.Pos `json:"pos"`
}

func (t *Token) String() string {
	return fmt.Sprintf("%05d %03d:%03d %-11s -> %#v",
		t.Pos.Byte, t.Pos.Line, t.Pos.Col,
		typeNames[t.Type],
		t.Value)
}
