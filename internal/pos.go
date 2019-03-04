package internal

import "fmt"

type Pos struct {
	Byte int `json"byte_no"`
	Line int `json:"line_no"`
	Col  int `json:"col_no"`
}

func (p Pos) String() string {
	return fmt.Sprintf("line %d col %d", p.Line, p.Col)
}
