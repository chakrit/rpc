package lexer

import (
	"fmt"

	"github.com/chakrit/rpc/internal"
)

type Token struct {
	Type  TokenType    `json:"type"`
	Value string       `json:"value"`
	Pos   internal.Pos `json:"pos"`
}

func (t *Token) String() string {
	return fmt.Sprintf("%05d %03d:%03d %-11s -> %#v",
		t.Pos.Byte, t.Pos.Line, t.Pos.Col,
		typeNames[t.Type],
		t.Value)
}
