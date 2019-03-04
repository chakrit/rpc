package lexer

func lexComment(c *lexer, r rune) lexFunc {
	c.Precond(IsCommentMarker(r), "expecting a / to start comments")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexComment_2
}

func lexComment_2(c *lexer, r rune) lexFunc {
	c.Precond(IsCommentMarker(r), "stray `/` need // for comments")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexComment_ContentStart
}

func lexComment_ContentStart(c *lexer, r rune) lexFunc {
	switch {
	case IsNewLine(r), IsEOF(r):
		return lexComment_End
	default:
		return lexComment_Content
	}
}

func lexComment_Content(c *lexer, r rune) lexFunc {
	c.Precond(!IsNewLine(r) && !IsEOF(r), "expecting valid comment characters")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexComment_ContentStart
}

func lexComment_End(c *lexer, r rune) lexFunc {
	c.Precond(IsNewLine(r) || IsEOF(r), "expecting termination of comment line")
	c.Emit(T_Comment, c.ClearBuffer())
	return lexStart
}
