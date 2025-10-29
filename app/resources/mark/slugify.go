package mark

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// Matches the TypeScript implementation in web/src/utils/slugify.ts

var (
	nonLetterNumberPattern = regexp.MustCompile(`[^\p{L}\p{M}\p{N}\-_]+`)
	multiHyphenPattern     = regexp.MustCompile(`-+`)
)

func Slugify(input string) string {
	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(input)

	// NFKC normalization
	normalized := norm.NFKC.String(trimmed)

	// Lowercase
	lowercased := strings.ToLower(normalized)

	// Replace non-letter/number chars with hyphens
	lettersReplaced := nonLetterNumberPattern.ReplaceAllString(lowercased, "-")

	// Collapse multiple hyphens
	collapsed := multiHyphenPattern.ReplaceAllString(lettersReplaced, "-")

	// Trim leading and trailing hyphens/underscores
	trimmedDividers := strings.Trim(collapsed, "-_")

	return trimmedDividers
}

// IsSlug checks if a string is already a valid slug by comparing it to its slugified version
func IsSlug(input string) bool {
	return Slugify(input) == input
}
