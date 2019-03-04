package lexer

func lexSpace(c *lexer, r rune) lexFunc {
	c.Precond(IsSpace(r), "expecting a non-newline whitespace character")
	return lexSpace_Start
}

func lexSpace_Start(c *lexer, r rune) lexFunc {
	switch {
	case IsSpace(r):
		return lexSpace_Content
	default:
		return lexSpace_End
	}
}

func lexSpace_Content(c *lexer, r rune) lexFunc {
	c.Precond(IsSpace(r), "expecting a non-newline whitespace character")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexSpace_Start
}

func lexSpace_End(c *lexer, r rune) lexFunc {
	c.Precond(!IsSpace(r), "expecting a non-whitespace character")
	c.Emit(T_Space, c.ClearBuffer())
	return lexStart
}
