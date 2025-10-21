package object

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestLocalStorerListMissingPrefixReturnsEmpty(t *testing.T) {
	t.Parallel()

	store := NewLocalStorer(config.Config{
		AssetStorageLocalPath: t.TempDir(),
	})

	got, err := store.List(context.Background(), "plugins")
	require.NoError(t, err)
	require.Empty(t, got)
}
