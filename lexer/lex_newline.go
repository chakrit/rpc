package lexer

func lexNewLine(c *lexer, r rune) lexFunc {
	switch {
	case IsCR(r):
		return lexNewLine_CR
	case IsLF(r):
		return lexNewLine_LF
	default:
		c.Precond(false, "expecting newline runes")
		return nil
	}
}

func lexNewLine_CR(c *lexer, r rune) lexFunc {
	c.Precond(IsCR(r), "expecting a CR")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexNewLine_CRLF
}

func lexNewLine_CRLF(c *lexer, r rune) lexFunc {
	if !IsLF(r) {
		return c.Fail("stray \\r, malformed line endings")
	}
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	c.Emit(T_EndOfLine, c.ClearBuffer())
	c.MarkNewLine()
	return lexStart
}

func lexNewLine_LF(c *lexer, r rune) lexFunc {
	c.Precond(IsLF(r), "expecting an LF")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	c.Emit(T_EndOfLine, c.ClearBuffer())
	c.MarkNewLine()
	return lexStart
}
