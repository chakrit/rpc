package lexer

import (
	"strings"
	"unicode"
)

var Keywords = []string{
	// meta
	"option",
	"include",
	"root", // unused, reserved as root ns identifier

	// definitions
	"namespace",
	"type",
	"rpc",

	// built-in types
	"string",
	"bool",
	"int",
	"long",
	"float",
	"double",
	"list",
	"map",
	"time",
}

var (
	IsEOF = func(r rune) bool { return r == 0 }
	IsCR  = func(r rune) bool { return r == '\r' }
	IsLF  = func(r rune) bool { return r == '\n' }

	IsNewLine = func(r rune) bool { return IsCR(r) || IsLF(r) }
	IsSpace   = func(r rune) bool { return !IsNewLine(r) && unicode.IsSpace(r) }
	IsBrace   = func(r rune) bool { return strings.ContainsRune("{}()<>", r) }

	IsCommentMarker      = func(r rune) bool { return r == '/' }
	IsStringMarker       = func(r rune) bool { return r == '"' }
	IsStringEscapeMarker = func(r rune) bool { return r == '\\' }

	IsDigit            = func(r rune) bool { return unicode.IsDigit(r) }
	IsDecimalSeparator = func(r rune) bool { return r == '.' }
	IsArgSeparator     = func(r rune) bool { return r == ',' }

	IsValidIdentFirstChar = func(r rune) bool { return unicode.IsLetter(r) || strings.ContainsRune("_-", r) }
	IsValidIdent          = func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune("_-", r) }
)

func IsKeyword(word string) bool {
	for _, str := range Keywords {
		if str == word {
			return true
		}
	}
	return false
}
