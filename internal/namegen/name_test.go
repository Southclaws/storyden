package namegen

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVocabulary(t *testing.T) {
	t.Parallel()
	permutations, longest := 1, 3
	word := regexp.MustCompile(`^[a-z]+$`)
	for _, words := range [][]string{moods[:], colours[:], places[:], characters[:]} {
		seen := map[string]bool{}
		maxLength := 0
		for _, w := range words {
			require.Regexp(t, word, w)
			require.False(t, seen[w], "duplicate word %q", w)
			seen[w] = true
			maxLength = max(maxLength, len(w))
		}
		longest += maxLength
		permutations *= len(words)
	}
	require.GreaterOrEqual(t, permutations, 16_000_000)
	require.LessOrEqual(t, longest, 30)
}

func TestGenerate(t *testing.T) {
	t.Parallel()
	pattern := regexp.MustCompile(`^[a-z]+-[a-z]+-[a-z]+-[a-z]+$`)
	for range 1000 {
		name := Generate()
		require.Regexp(t, pattern, name)
		require.LessOrEqual(t, len(name), 30)
	}
}
