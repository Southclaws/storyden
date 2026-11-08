// Package case_convert provides casing conversions for identifiers.
// Extend this package when new casing operations are needed.
package case_convert

import (
	"strings"
	"unicode"
)

// ToCamel joins delimiter-separated words in CamelCase: HTTP_SERVER becomes
// HttpServer. Each word is lowercased with its first character uppercased.
// Unicode letters and digits are preserved; everything else separates words.
// Existing case transitions and digits do not introduce word boundaries.
func ToCamel(value string) string {
	words := strings.FieldsFunc(value, func(character rune) bool {
		return !unicode.IsLetter(character) && !unicode.IsDigit(character)
	})

	var result strings.Builder
	result.Grow(len(value))

	for _, word := range words {
		for index, character := range word {
			if index == 0 {
				result.WriteRune(unicode.ToUpper(character))
			} else {
				result.WriteRune(unicode.ToLower(character))
			}
		}
	}

	return result.String()
}
