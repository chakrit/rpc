package parser

import "github.com/chakrit/rpc/lexer"

func (p *parser) parseOption() (string, interface{}, error) {
	t := p.Peek()
	p.Precond(t.Value == "option", "expecting `option` keyword")

	_, tkey := p.Consume()
	if tkey.Type != lexer.T_Identifier {
		return "", nil, p.Fail("option name expected")
	}

	_, tval := p.Consume()
	if (tval.Type & (lexer.T_NumberValue + lexer.T_StringValue)) == 0 {
		return "", nil, p.Fail("option value literal expected")
	}

	key, val := tkey.Value, tval.Value
	p.Consume()

	return key, val, nil
}
