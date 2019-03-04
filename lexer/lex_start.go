package lexer

func lexStart(c *lexer, r rune) lexFunc {
	switch {
	case IsEOF(r):
		return nil
	case IsBrace(r):
		return lexBrace
	case IsArgSeparator(r):
		return lexSeparator
	case IsCommentMarker(r):
		return lexComment
	case IsStringMarker(r):
		return lexString
	case IsNewLine(r):
		return lexNewLine
	case IsSpace(r):
		return lexSpace
	case IsDigit(r):
		return lexNumber
	case IsValidIdentFirstChar(r):
		return lexIdentifier
	default:
		return c.Fail("bad unicode sequence `" + string(r) + "`")
	}
}

func lexBrace(c *lexer, r rune) lexFunc {
	c.Precond(IsBrace(r), "expecting scope begin/end char")
	if !c.Consume() { // valid rune by this point
		return nil
	}

	typ, ok := braceMappings[r]
	if !ok {
		return c.Fail("invalid brace `" + string(r) + "`")
	}

	c.Emit(typ, string(r))
	return lexStart
}

func lexSeparator(c *lexer, r rune) lexFunc {
	c.Precond(IsArgSeparator(r), "expecting separator char")
	if !c.Consume() {
		return nil
	}

	c.Emit(T_ArgListSep, string(r))
	return lexStart
}
