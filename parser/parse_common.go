package parser

import (
	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/spec"
)

func (p *parser) parseBlockStart(scope string) (string, error) {
	t := p.Peek()
	p.Precond(t.Value == scope, "expecting `"+scope+"` keyword")

	_, ident := p.Consume()
	if ident.Type != lexer.T_Identifier {
		return "", p.Fail(scope + " name expected")
	}

	_, open := p.Consume()
	if open.Type != lexer.T_BlockStart {
		return "", p.Fail("opening brace for " + scope + " `{` expected")
	}

	p.Consume()
	return ident.Value, nil
}

func (p *parser) parseTypeRef(scope string) (*spec.TypeRef, error) {
	t := p.Peek()
	p.Precond(t.Type&(lexer.T_Identifier|lexer.T_Keyword) > 0, "expecting identifier or keyword")

	ref := &spec.TypeRef{Name: t.Value}
	p.Consume()

	t = p.Peek()
	if t.Type == lexer.T_TypeArgListStart {
		if err := p.parseTypeRef_ArgList(scope, ref); err != nil {
			return nil, err
		}
	}

	return ref, nil
}

func (p *parser) parseTypeRef_ArgList(scope string, ref *spec.TypeRef) error {
	open := p.Peek()
	p.Precond(open.Type == lexer.T_TypeArgListStart, "expecting `<` to start type arg list")

	ref.Arguments = nil
	p.Consume()

	for {
		t := p.Peek()
		if t.Type&(lexer.T_Identifier|lexer.T_Keyword) == 0 {
			return p.Fail("type argument expected")
		}

		if arg, err := p.parseTypeRef(scope); err != nil {
			return err
		} else {
			ref.Arguments = append(ref.Arguments, arg)
		}

		t = p.Peek()
		switch t.Type {
		case lexer.T_ArgListSep:
			p.Consume()
		case lexer.T_TypeArgListEnd:
			p.Consume()
			return nil
		default:
			return p.Fail("more type arguments with `,` or closing angle bracket `>` expected")
		}
	}
}
