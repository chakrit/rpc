package internal

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gobuffalo/flect"
)

// TODO: Make the character ranges unicode-aware.
var (
	urlPart = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9.]*([^a-zA-Z0-9.]+)?`)
)

// InflectPascal takes a word and produce a pascal-cased version of it.
//
// Example: go_lang -> GoLanguage
func InflectPascal(s string) string { return flect.Pascalize(s) }

// InflectSnake takes a word and produces a snake-cased version of it (separate words with `_`)
//
// Example: GoLanguage -> go_language
func InflectSnake(s string) string { return flect.Underscore(s) }

// MangleID expects an url-like string and generates a mangled identifier that has a high
// chance to be valid and unique in most scopes.
//
// Example: github.com/chakrit/rpc -> gc_rpc
func MangleID(s string, inserts ...string) string {
	builder := strings.Builder{}
	parts := urlPart.FindAllString(s, -1)
	for _, part := range parts[:len(parts)-1] {
		r, _ := utf8.DecodeRuneInString(part)
		builder.WriteRune(r)
	}
	for _, insert := range inserts {
		builder.WriteString(insert)
	}

	builder.WriteRune('_')
	builder.WriteString(parts[len(parts)-1])
	return builder.String()
}
