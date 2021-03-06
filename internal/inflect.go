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

// InflectDash takes a word and produces a dasherized version of it (separate words with `-`)
func InflectDash(s string) string { return flect.Dasherize(s) }

// InflectTitle capitalize words and put spaces between them.
func InflectTitle(s string) string { return flect.Titleize(s) }
