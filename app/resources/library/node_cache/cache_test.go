package node_cache

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"
)

func TestRevisionRejectsValidatorStoredByEarlierRead(t *testing.T) {
	ctx := context.Background()
	cache := New(slog.New(slog.NewTextHandler(io.Discard, nil)), cachetest.New())
	key := "child-page"

	etag, err := cache.Prepare(ctx, key)
	require.NoError(t, err)
	require.NoError(t, cache.Store(ctx, key, etag))

	query := cachecontrol.NewQuery(opt.New(etag.String()))
	_, notModified := cache.Check(ctx, query, key)
	require.True(t, notModified)

	require.NoError(t, cache.InvalidateBeforeWrite(ctx))
	require.NoError(t, cache.Store(ctx, key, etag))

	_, notModified = cache.Check(ctx, query, key)
	assert.False(t, notModified)
}

func TestCanonicalKey(t *testing.T) {
	id := xid.New()

	tests := map[string]struct {
		key  mark.Queryable
		want string
	}{
		"slug": {
			key:  mark.NewQueryKey("node-slug"),
			want: "node-slug",
		},
		"id": {
			key:  mark.NewQueryKeyID(id),
			want: id.String(),
		},
		"id and slug": {
			key:  mark.NewQueryKey(id.String() + "-node-slug"),
			want: id.String(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, CanonicalKey(tt.key))
		})
	}
}
