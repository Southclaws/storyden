package node_cache

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidateManyRefreshesCurrentAndRetiresOldAliases(t *testing.T) {
	ctx := context.Background()
	store := newRecordingStore()
	now := time.Date(2026, time.September, 26, 12, 0, 0, 0, time.UTC)
	cache := &Cache{store: store, clock: func() time.Time { return now }}
	first, second := xid.New(), xid.New()
	store.values[cachePrefix+"retired"] = now.Add(-time.Hour).Format(storeTimeFmt)

	require.NoError(t, cache.InvalidateMany(ctx, []Invalidation{
		{ID: first, Slug: "current", PreviousSlug: "retired"},
		{ID: second, Slug: "retained", PreviousSlug: "current"},
	}))

	assert.Equal(t, 1, store.batchCalls)
	assert.Nil(t, cache.lastModified(ctx, "retired"))
	for _, key := range []string{first.String(), second.String(), "current", "retained"} {
		require.NotNil(t, cache.lastModified(ctx, key))
		assert.Equal(t, now, *cache.lastModified(ctx, key))
	}
}
