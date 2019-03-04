package lexer

func lexString(c *lexer, r rune) lexFunc {
	c.Precond(IsStringMarker(r), "expecting a quotation mark to start string")
	if !c.Consume() {
		return nil
	}

	return lexString_ContentStart
}

func lexString_ContentStart(c *lexer, r rune) lexFunc {
	switch {
	case IsEOF(r):
		return c.Fail("unterminated string at end of file")
	case IsStringMarker(r):
		return lexString_End
	case IsStringEscapeMarker(r):
		return lexString_EscapedChar
	default:
		return lexString_Content
	}
}

func lexString_Content(c *lexer, r rune) lexFunc {
	c.Precond(!IsStringMarker(r) && !IsStringEscapeMarker(r), "expecting non-escaped string char")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexString_ContentStart
}

func lexString_EscapedChar(c *lexer, r rune) lexFunc {
	if !IsStringEscapeMarker(r) {
		return c.Fail("unterminated string at end of file")
	}
	if !c.Consume() {
		return nil
	}

	return lexString_EscapedChar2
}

func lexString_EscapedChar2(c *lexer, r rune) lexFunc {
	if IsEOF(r) {
		return c.Fail("unterminated string at end of file")
	}
	if !c.Consume() {
		return nil
	}

	switch r {
	case 'r':
		c.AppendBuffer('\r')
	case 'n':
		c.AppendBuffer('\n')
	case 't':
		c.AppendBuffer('\t')
	default: // other chars have no meaning, they're added as-is ['"abc123\/]
		c.AppendBuffer(r)
	}

	return lexString_ContentStart
}

func lexString_End(c *lexer, r rune) lexFunc {
	c.Precond(IsStringMarker(r), "expecting string termination char")
	if !c.Consume() {
		return nil
	}

	c.Emit(T_StringValue, c.ClearBuffer())
	return lexStart
}
