package token

import (
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratedTokensAreUnpredictable(t *testing.T) {
	t.Parallel()

	const runs = 512

	seen := map[string]struct{}{}
	prefixes := map[string]struct{}{}

	for range runs {
		tok, err := Generate()
		require.NoError(t, err)

		seen[tok.String()] = struct{}{}
		prefixes[tok.String()[:8]] = struct{}{}
	}

	assert.Len(t, seen, runs, "every secret must be distinct")
	assert.Len(t, prefixes, runs, "no shared structure, an xid would collide here on machine and pid bytes")
}

func TestTokenIsNotAnXID(t *testing.T) {
	t.Parallel()

	tok, err := Generate()
	require.NoError(t, err)

	_, err = xid.FromString(tok.String())
	assert.Error(t, err, "a session secret must not be parseable as the sequential id it used to be")
}

func TestFromStringRejectsMalformedSecrets(t *testing.T) {
	t.Parallel()

	tok, err := Generate()
	require.NoError(t, err)

	_, err = FromString(tok.String())
	assert.NoError(t, err)

	for name, input := range map[string]string{
		"empty":      "",
		"short":      tok.String()[:10],
		"long":       tok.String() + "a",
		"an xid":     xid.New().String(),
		"non base64": strings.Repeat("!", secretLength),
	} {
		_, err := FromString(input)
		assert.Error(t, err, name)
	}
}

func TestHashIsStableAndHidesTheSecret(t *testing.T) {
	t.Parallel()

	tok, err := Generate()
	require.NoError(t, err)

	assert.Equal(t, tok.Hash(), tok.Hash())
	assert.NotContains(t, tok.Hash(), tok.String())
	assert.Len(t, tok.Hash(), 64)

	other, err := Generate()
	require.NoError(t, err)
	assert.NotEqual(t, tok.Hash(), other.Hash())
}

func TestSerialisedSessionCarriesNoSecret(t *testing.T) {
	t.Parallel()

	tok, err := Generate()
	require.NoError(t, err)

	payload, err := Session{TokenHash: tok.Hash()}.Serialise()
	require.NoError(t, err)

	assert.NotContains(t, string(payload), tok.String(), "the cache must never hold a replayable secret")
	assert.Contains(t, string(payload), tok.Hash())
}
