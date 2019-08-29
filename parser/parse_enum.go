package parser

import (
	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/spec"
)

func (p *parser) parseEnum() (*spec.Enum, error) {
	name, err := p.parseBlockStart("enum")
	if err != nil {
		return nil, err
	}

	enum := &spec.Enum{Name: name}
	if err := p.parseEnum_Members(enum); err != nil {
		return nil, err
	}

	closing, _ := p.Consume()
	if closing.Type != lexer.T_BlockEnd {
		return nil, p.Fail("closing bracket for enum{} expected")
	}

	return enum, nil
}

func (p *parser) parseEnum_Members(enum *spec.Enum) error {
	for {
		t := p.Peek()
		switch t.Type {
		case lexer.T_Keyword, lexer.T_Identifier:
			enum.Members = append(enum.Members, t.Value)
			p.Consume()
		case lexer.T_BlockEnd:
			return nil
		case lexer.T_EndOfFile:
			return p.Fail("missing closing brace for enum{}")
		default:
			return p.Fail("enum member definition expected")
		}
	}
}
