package parser

import (
	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/spec"
)

func (p *parser) parseType() (*spec.Type, error) {
	name, err := p.parseBlockStart("type")
	if err != nil {
		return nil, err
	}

	typ := &spec.Type{Name: name}
	if err := p.parseType_Content(typ); err != nil {
		return nil, err
	}

	closing, _ := p.Consume()
	if closing.Type != lexer.T_BlockEnd {
		return nil, p.Fail("closing bracket for type{} expected")
	}

	return typ, nil
}

func (p *parser) parseType_Content(typ *spec.Type) error {
	for {
		t := p.Peek()
		switch t.Type {
		case lexer.T_Keyword, lexer.T_Identifier:
			// continue
		case lexer.T_BlockEnd:
			return nil
		case lexer.T_EndOfFile:
			return p.Fail("missing closing brace for type{}")
		default:
			return p.Fail("property definition expected")
		}

		typeref, err := p.parseTypeRef("type")
		if err != nil {
			return err
		}

		ident := p.Peek()
		if ident.Type&(lexer.T_Identifier|lexer.T_Keyword) == 0 {
			return p.Fail("property name expected")
		}

		prop := &spec.Property{
			Name: ident.Value,
			Type: typeref,
		}
		_, isNew := typ.Properties.AddIfNew(prop)
		if !isNew {
			return p.Fail("duplicate declaration for property `" + prop.Name + "`")
		}

		typ.Properties[prop.Name] = prop
		p.Consume()
	}
}
