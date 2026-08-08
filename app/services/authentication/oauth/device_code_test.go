package oauth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateUserCode(t *testing.T) {
	require.Len(t, userCodeAlphabet, 32)
	for _, excluded := range "ILOU-_" {
		assert.NotContains(t, userCodeAlphabet, string(excluded))
	}
	assert.Contains(t, userCodeAlphabet, "0")
	assert.Contains(t, userCodeAlphabet, "1")

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

		// Crockford omits ambiguous letters and accepts them as input aliases.
		assert.NotContains(t, raw, "O")
		assert.NotContains(t, raw, "I")
		assert.NotContains(t, raw, "L")
		assert.NotContains(t, raw, "U")
	}
}

func TestNormalizeCode(t *testing.T) {
	tests := map[string]string{
		"ABCD-EFGH":   "ABCDEFGH",
		" abcd-efgh ": "ABCDEFGH",
		"O234-5678":   "02345678",
		"I234-L678":   "12341678",
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, expected, normalizeCode(input))
		})
	}
}

func TestNormalizeCodeRoundTripsGeneratedCodes(t *testing.T) {
	for range 1000 {
		userCode, err := generateUserCode()
		require.NoError(t, err)

		normalized := normalizeCode(userCode)

		assert.Equal(t, strings.ReplaceAll(userCode, "-", ""), normalized)
		assert.Len(t, normalized, 8)
	}
}

func TestParseUserCode(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
		ok    bool
	}{
		"canonical":           {input: "ABCD-EFGH", want: "ABCDEFGH", ok: true},
		"lowercase":           {input: " abcd-efgh ", want: "ABCDEFGH", ok: true},
		"crockford aliases":   {input: "O234-I67L", want: "02341671", ok: true},
		"separators ignored":  {input: "AB-CD-EF-GH", want: "ABCDEFGH", ok: true},
		"missing":             {input: "", ok: false},
		"short":               {input: "ABCD-EFG", ok: false},
		"long":                {input: "ABCD-EFGH-J", ok: false},
		"invalid punctuation": {input: "ABCD_EFG", ok: false},
		"excluded U":          {input: "ABCD-UEFG", ok: false},
		"non ASCII":           {input: "ABCD-ÉFGH", ok: false},
		"oversized":           {input: strings.Repeat("A", maxUserCodeInputLength+1), ok: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := parseUserCode(test.input)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.want, got)
		})
	}
}
