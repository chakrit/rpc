package lexer

import (
	"strings"
	"unicode"
)

func IsEOF(r rune) bool { return r == 0 }
func IsCR(r rune) bool  { return r == '\r' }
func IsLF(r rune) bool  { return r == '\n' }

func IsNewLine(r rune) bool { return IsCR(r) || IsLF(r) }
func IsSpace(r rune) bool   { return !IsNewLine(r) && unicode.IsSpace(r) }
func IsBrace(r rune) bool   { return strings.ContainsRune("{}()<>", r) }

func IsCommentMarker(r rune) bool      { return r == '/' }
func IsStringMarker(r rune) bool       { return r == '"' }
func IsStringEscapeMarker(r rune) bool { return r == '\\' }

func IsDigit(r rune) bool            { return unicode.IsDigit(r) }
func IsDecimalSeparator(r rune) bool { return r == '.' }
func IsArgSeparator(r rune) bool     { return r == ',' }

func IsValidIdentFirstChar(r rune) bool { return unicode.IsLetter(r) || strings.ContainsRune("_-", r) }
func IsValidIdent(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune("_-", r)
}
