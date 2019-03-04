package lexer

import (
	"bufio"
	"io"

	"github.com/chakrit/rpc/internal"
	"github.com/pkg/errors"
)

type lexFunc func(*lexer, rune) lexFunc

type Options struct {
	Input       io.Reader
	Logger      internal.Logger
	IgnoreTypes int
}

type lexer struct {
	logger  internal.Logger
	reader  *bufio.Reader
	ignores int
	pos     internal.Pos

	buffer string
	head   rune
	err    error

	state  lexFunc
	tokens []*Token
}

func Lex(opts Options) ([]*Token, error) {
	ctx, err := newLexer(opts, lexStart)
	if err != nil {
		return nil, errors.Wrap(err, ctx.pos.String())
	}

	for ctx.Step() {
	}

	if ctx.Err() != nil {
		return nil, errors.Wrapf(ctx.Err(), ctx.pos.String())
	} else {
		ctx.Emit(T_EndOfFile, "")
		return ctx.tokens, nil
	}
}

func newLexer(opts Options, startState lexFunc) (*lexer, error) {
	ctx := &lexer{
		logger:  opts.Logger,
		reader:  bufio.NewReader(opts.Input),
		ignores: opts.IgnoreTypes,
		state:   startState,
	}

	// get initial peek rune, don't Consume() since we don't want to increment position
	if r, _, err := ctx.reader.ReadRune(); err != nil && err != io.EOF {
		return nil, err
	} else {
		ctx.head = r
	}

	return ctx, nil
}

func (c *lexer) Precond(cond bool, msg string) {
	if !cond {
		err := errors.New("precondition failure: " + msg)
		c.logger.Fatalp(c.pos, err)
	}
}

// error handling
func (c *lexer) Err() error { return c.err }
func (c *lexer) Fail(msg string) lexFunc {
	c.err = errors.New(msg)
	return nil
}

// buffer handling
func (c *lexer) MarkNewLine() {
	c.pos = internal.Pos{
		Byte: c.pos.Byte,
		Line: c.pos.Line + 1,
		Col:  0,
	}
}

func (c *lexer) AppendBuffer(r rune) { c.buffer = c.buffer + string(r) }
func (c *lexer) ClearBuffer() string {
	str := c.buffer
	c.buffer = ""
	return str
}

// state handling
func (c *lexer) Consume() bool {
	r, n, err := c.reader.ReadRune()
	if err != nil && err != io.EOF {
		c.err = errors.Wrap(err, "read failure")
		return false
	}

	c.pos = internal.Pos{
		Byte: c.pos.Byte + n,
		Col:  c.pos.Col + 1,
		Line: c.pos.Line,
	}
	if err == io.EOF {
		c.pos.Byte += 1 // virtual byte num for the EOF mark
		c.head = 0
	} else {
		c.head = r
	}
	return true
}

func (c *lexer) Step() bool {
	c.state = c.state(c, c.head)
	return c.state != nil && c.err == nil
}

func (c *lexer) Emit(typ int, value string) {
	if c.ignores&typ > 0 {
		return
	}

	c.tokens = append(c.tokens, &Token{
		Type:  typ,
		Value: value,
		Pos:   c.pos,
	})
}
