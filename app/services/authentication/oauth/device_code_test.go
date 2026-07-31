package oauth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateUserCode(t *testing.T) {
	for i := 0; i < 1000; i++ {
		userCode, err := generateUserCode()
		require.NoError(t, err)

		assert.Len(t, userCode, 9)
		assert.Equal(t, "-", string(userCode[4]))

		raw := strings.ReplaceAll(userCode, "-", "")
		assert.Len(t, raw, 8)

		for _, c := range raw {
			assert.Contains(t, userCodeAlphabet, string(c), "character %q must come from userCodeAlphabet", c)
		}

		// The 0/O/1/I/L/U confusables must never appear in generated codes,
		// and the hyphen separator must be the only one in the display code.
		assert.NotContains(t, raw, "O")
		assert.NotContains(t, raw, "I")
		assert.NotContains(t, raw, "L")
		assert.NotContains(t, raw, "U")
	}
}

func TestNormalizeCodeRoundTripsGeneratedCodes(t *testing.T) {
	userCode, err := generateUserCode()
	require.NoError(t, err)

	normalized := normalizeCode(userCode)

	assert.Equal(t, strings.ReplaceAll(userCode, "-", ""), normalized)
	assert.Len(t, normalized, 8)
}
