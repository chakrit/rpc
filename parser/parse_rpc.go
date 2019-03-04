package parser

import (
	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/spec"
)

func (p *parser) parseRPC() (*spec.RPC, error) {
	t := p.Peek()
	p.Precond(t.Value == "rpc", "expecting `rpc` context")

	_, ident := p.Consume()
	if ident.Type != lexer.T_Identifier {
		return nil, p.Fail("rpc name expected")
	}

	_, open := p.Consume()
	if open.Type != lexer.T_ArgListStart {
		return nil, p.Fail("start of argument list `(` expected")
	}

	rpc := &spec.RPC{Name: ident.Value}
	p.Consume()

	if err := p.parseRPC_InputArgs(rpc); err != nil {
		return nil, err
	} else if err := p.parseRPC_OutputArgs(rpc); err != nil {
		return nil, err
	}

	return rpc, nil
}

func (p *parser) parseRPC_InputArgs(rpc *spec.RPC) error {
	rpc.InputTypes = nil

	for {
		t := p.Peek()
		switch t.Type {
		case lexer.T_Identifier, lexer.T_Keyword:
			// continue
		case lexer.T_ArgListEnd:
			p.Consume()
			return nil
		default:
			return p.Fail("argument types expected")
		}

		if ref, err := p.parseTypeRef("rpc"); err != nil {
			return err
		} else {
			rpc.InputTypes = append(rpc.InputTypes, ref)
		}

		t = p.Peek()
		switch t.Type {
		case lexer.T_ArgListSep:
			p.Consume()
		case lexer.T_ArgListEnd:
			p.Consume()
			return nil
		default:
			return p.Fail("more input types with `,` or closing bracket `)` expected")
		}
	}
}

func (p *parser) parseRPC_OutputArgs(r *spec.RPC) error {
	r.OutputTypes = nil
	if ref, err := p.parseTypeRef("rpc"); err != nil {
		return err
	} else {
		r.OutputTypes = []*spec.TypeRef{ref}
		return nil
	}
}
