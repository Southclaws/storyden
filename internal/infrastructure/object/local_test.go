package object

import (
	"context"
	"strings"
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

func TestLocalStorerListTrailingSlash(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := NewLocalStorer(config.Config{AssetStorageLocalPath: t.TempDir()})
	for _, key := range []string{"avatar/member", "avatar/nested/member"} {
		require.NoError(t, store.Write(ctx, key, strings.NewReader("image"), 5))
	}

	for _, prefix := range []string{"avatar", "avatar/"} {
		got, err := store.List(ctx, prefix)
		require.NoError(t, err)
		require.ElementsMatch(t, []string{"member", "nested"}, got)
	}
}
