package node_mutate

import (
	"io"
	"log/slog"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func TestContentWithStableBlocksPreservesBlockIDsForImageOnlyContent(t *testing.T) {
	previousParsed, err := datagraph.NewRichText(`<img src="https://example.com/image.jpg" alt="Before">`)
	require.NoError(t, err)
	previous, err := datagraph.NewRichTextWithNewBlocks(previousParsed)
	require.NoError(t, err)
	require.False(t, previous.Content.IsEmpty())

	next, err := datagraph.NewRichText(`<img src="https://example.com/image.jpg" alt="After">`)
	require.NoError(t, err)

	updated := contentWithStableBlocks(
		opt.New(previous.Content),
		next,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	assert.Equal(t, firstBlockIDOfType(t, previous.Content, "img"), firstBlockIDOfType(t, updated, "img"))
}

func firstBlockIDOfType(t *testing.T, content datagraph.Content, typ string) string {
	t.Helper()
	for _, block := range content.Blocks() {
		if block.Type == typ {
			return block.ID
		}
	}
	t.Fatalf("no %s block found", typ)
	return ""
}
