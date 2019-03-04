package lexer

func lexNumber(c *lexer, r rune) lexFunc {
	c.Precond(IsDigit(r), "expecting a digit")
	return lexNumber_Start
}

func lexNumber_Start(c *lexer, r rune) lexFunc {
	switch {
	case IsDigit(r):
		return lexNumber_Int
	case IsDecimalSeparator(r):
		return lexNumber_Separator
	default:
		return lexNumber_End
	}
}

func lexNumber_Int(c *lexer, r rune) lexFunc {
	c.Precond(IsDigit(r), "expecting a digit")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexNumber_Start
}

func lexNumber_Separator(c *lexer, r rune) lexFunc {
	c.Precond(IsDecimalSeparator(r), "expecting a dot")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexNumber_FractionStart
}

func lexNumber_FractionStart(c *lexer, r rune) lexFunc {
	switch {
	case IsDigit(r):
		return lexNumber_Fraction
	case IsDecimalSeparator(r):
		return c.Fail("multiple decimal separator in number literal")
	default:
		return lexNumber_End
	}
}

func lexNumber_Fraction(c *lexer, r rune) lexFunc {
	c.Precond(IsDigit(r), "expecting a digit")
	if !c.Consume() {
		return nil
	}

	c.AppendBuffer(r)
	return lexNumber_FractionStart
}

func lexNumber_End(c *lexer, r rune) lexFunc {
	c.Precond(!IsDigit(r) && !IsDecimalSeparator(r), "expecting end of number sequence")
	c.Emit(T_NumberValue, c.ClearBuffer())
	return lexStart
}
