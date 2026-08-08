package local

import (
	"context"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/stretchr/testify/require"
)

func TestWriteMethodsReturnRejectedWrites(t *testing.T) {
	ctx := context.Background()

	t.Run("set", func(t *testing.T) {
		cache := newTestLocalCache(t)
		cache.cache.Close()

		require.ErrorIs(t, cache.Set(ctx, "key", "value", time.Minute), errWriteRejected)
	})

	t.Run("set many", func(t *testing.T) {
		cache := newTestLocalCache(t)
		cache.cache.Close()

		require.ErrorIs(t, cache.SetMany(ctx, map[string]string{"key": "value"}, time.Minute), errWriteRejected)
	})

	t.Run("hash set", func(t *testing.T) {
		cache := newTestLocalCache(t)
		cache.cache.Close()

		_, err := cache.HIncrBy(ctx, "key", "field", 1)
		require.ErrorIs(t, err, errWriteRejected)
	})

	t.Run("expire", func(t *testing.T) {
		cache := newTestLocalCache(t)
		require.NoError(t, cache.Set(ctx, "key", "value", time.Minute))

		require.ErrorIs(t, cache.Expire(ctx, "key", -time.Second), errWriteRejected)
	})
}

func newTestLocalCache(t *testing.T) *LocalCache {
	t.Helper()

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 100,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	require.NoError(t, err)
	t.Cleanup(cache.Close)

	return &LocalCache{cache: cache}
}
