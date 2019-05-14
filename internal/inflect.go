package internal

import (
	"github.com/gobuffalo/flect"
)

// InflectPascal takes a word and produce a pascal-cased version of it.
//
// Example: go_lang -> GoLanguage
func InflectPascal(s string) string { return flect.Pascalize(s) }

// InflectSnake takes a word and produces a snake-cased version of it (separate words with `_`)
//
// Example: GoLanguage -> go_language
func InflectSnake(s string) string { return flect.Underscore(s) }
