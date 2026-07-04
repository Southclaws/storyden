package node_versioning

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/deletable"
)

func TestApplyDraftPartialPreservesBlockIDsForImageOnlyContent(t *testing.T) {
	previousParsed, err := datagraph.NewRichText(`<img src="https://example.com/image.jpg" alt="Before">`)
	require.NoError(t, err)
	previous, err := datagraph.NewRichTextWithNewBlocks(previousParsed)
	require.NoError(t, err)
	require.False(t, previous.Content.IsEmpty())

	next, err := datagraph.NewRichText(`<img src="https://example.com/image.jpg" alt="After">`)
	require.NoError(t, err)

	snapshot := draftSnapshot{
		Content: opt.New(previous.Content),
	}
	err = applyDraftPartial(&snapshot, DraftPartial{
		Content: deletable.Skip(opt.New(next)),
	})
	require.NoError(t, err)

	updated, ok := snapshot.Content.Get()
	require.True(t, ok)
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
