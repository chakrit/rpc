package parser

import (
	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/spec"
)

func (p *parser) parseRoot() (*spec.Namespace, error) {
	ns := &spec.Namespace{Name: "root"}
	if err := p.parseNamespace_Content(ns); err != nil {
		return nil, err
	} else {
		return ns, nil
	}
}

func (p *parser) parseNamespace() (*spec.Namespace, error) {
	name, err := p.parseBlockStart("namespace")
	if err != nil {
		return nil, err
	}

	ns := &spec.Namespace{Name: name}
	if err := p.parseNamespace_Content(ns); err != nil {
		return nil, err
	}

	closing, _ := p.Consume()
	if closing.Type != lexer.T_BlockEnd {
		return nil, p.Fail("missing closing bracket `}` for namespace{}")
	}

	return ns, nil
}

func (p *parser) parseNamespace_Content(ns *spec.Namespace) error {
	for {
		t := p.Peek()
		switch t.Type {
		case lexer.T_Keyword:
			// continue to keyword switch
		case lexer.T_BlockEnd, lexer.T_EndOfFile:
			return nil
		default:
			return p.Fail("valid definition keyword expected")
		}

		switch t.Value {
		case "namespace":
			if child, err := p.parseNamespace(); err != nil {
				return err
			} else {
				ns.Children.Add(child)
			}

		case "option":
			key, value, err := p.parseOption()
			if err != nil {
				return err
			}
			if ns.Options == nil {
				ns.Options = map[string]interface{}{}
			}
			ns.Options[key] = value

		case "type":
			if typ, err := p.parseType(); err != nil {
				return err
			} else {
				ns.Types.Add(typ)
			}

		case "rpc":
			if r, err := p.parseRPC(); err != nil {
				return err
			} else if _, isNew := ns.RPCs.AddIfNew(r); !isNew {
				return p.Fail("duplicate definition for RPC `" + r.Name + "`")
			} // else ok, rpc added

		default:
			return p.Fail("unrecognized keyword: `" + t.Value + "`")
		}
	}
}
