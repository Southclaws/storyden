package sse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestStreamOffsetRoundTrip(t *testing.T) {
	for _, offset := range []uint64{0, 1, 42, 9999999999999999} {
		parsed, err := parseStreamOffset(formatStreamOffset(offset))
		require.NoError(t, err)
		assert.Equal(t, offset, parsed)
	}
	for _, invalid := range []string{"nope", "1", "0000000000000001_0000000000000001"} {
		_, err := parseStreamOffset(invalid)
		assert.ErrorIs(t, err, errInvalidStreamOffset)
	}
}

func TestJSONBatch(t *testing.T) {
	first := openapi.StreamPart{}
	require.NoError(t, first.FromTextStartPart(openapi.TextStartPart{Id: "text"}))
	second := openapi.StreamPart{}
	require.NoError(t, second.FromTextDeltaPart(openapi.TextDeltaPart{Id: "text", Delta: "hello"}))

	assert.JSONEq(t, `[{"type":"text-start","id":"text"},{"type":"text-delta","id":"text","delta":"hello"}]`, string(jsonBatch([]openapi.StreamPart{first, second})))
}
