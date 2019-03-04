package lexer

func lexIdentifier(c *lexer, r rune) lexFunc {
	c.Precond(IsValidIdentFirstChar(r), "expecting start of identifier")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexIdentifier_NextCharStart
}

func lexIdentifier_NextCharStart(c *lexer, r rune) lexFunc {
	switch {
	case IsValidIdent(r):
		return lexIdentifier_NextChar
	default:
		return lexIdentifier_End
	}
}

func lexIdentifier_NextChar(c *lexer, r rune) lexFunc {
	c.Precond(IsValidIdent(r), "expecting valid identifier character")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexIdentifier_NextCharStart
}

func lexIdentifier_End(c *lexer, r rune) lexFunc {
	c.Precond(!IsValidIdent(r), "expecting an end of identifier")

	word := c.ClearBuffer()
	if IsKeyword(word) {
		c.Emit(T_Keyword, word)
	} else {
		c.Emit(T_Identifier, word)
	}
	return lexStart
}
