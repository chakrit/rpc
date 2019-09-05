package parser

import (
	"fmt"
	"io"

	"github.com/chakrit/rpc/internal"

	"errors"

	"github.com/chakrit/rpc/lexer"
	"github.com/chakrit/rpc/spec"
)

type Options struct {
	Input  io.Reader
	Logger internal.Logger
}

type parser struct {
	logger internal.Logger
	tokens []*lexer.Token
	debug  bool
	pos    int
}

func Parse(opts Options) (*spec.Namespace, error) {
	tokens, err := lexer.Lex(lexer.Options{
		Input:       opts.Input,
		Logger:      opts.Logger,
		IgnoreTypes: lexer.T_Space + lexer.T_EndOfLine + lexer.T_Comment,
	})
	if err != nil {
		return nil, err
	}

	p := &parser{
		logger: opts.Logger,
		tokens: tokens,
	}
	ns, err := p.parseRoot()
	if err != nil {
		if token := p.Peek(); token != nil {
			return nil, fmt.Errorf(token.Pos.String()+": %w", err)
		} else {
			return nil, fmt.Errorf("parse failure: %w", err)
		}
	}

	return ns, nil
}

func (p *parser) Precond(cond bool, msg string) {
	if !cond {
		err := errors.New("precondition failure: " + msg)
		if token := p.Peek(); token != nil {
			p.logger.Fatalp(token.Pos, err)
		} else {
			p.logger.Fatal(err)
		}
	}
}

func (p *parser) Fail(msg string) error {
	t := p.Peek()
	return fmt.Errorf("near `%s`: %s", t.Value, msg)
}

func (p *parser) Peek() *lexer.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	} else {
		return nil
	}
}

func (p *parser) Consume() (*lexer.Token, *lexer.Token) {
	t := p.Peek()
	if t == nil {
		return nil, nil
	} else {
		p.pos += 1
		return t, p.Peek()
	}
}
