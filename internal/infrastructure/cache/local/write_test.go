package local

import (
	"context"
	"fmt"
	"sync"
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

	t.Run("set_if_absent", func(t *testing.T) {
		cache := newTestLocalCache(t)
		cache.cache.Close()

		inserted, err := cache.SetIfAbsent(ctx, "key", "value", time.Minute)
		require.ErrorIs(t, err, errWriteRejected)
		require.False(t, inserted)
	})

	t.Run("set_if_absent_admission_rejected", func(t *testing.T) {
		cache := newTestLocalCache(t)
		cache.cache.UpdateMaxCost(1)

		inserted, err := cache.SetIfAbsent(ctx, "key", "value", time.Minute)
		require.ErrorIs(t, err, errWriteRejected)
		require.False(t, inserted)
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

func TestWritesOnDifferentStripesProceedIndependently(t *testing.T) {
	cache := newTestLocalCache(t)
	blocked := cache.lockIndex("blocked")
	key := "other"
	for i := 0; cache.lockIndex(key) == blocked; i++ {
		key = fmt.Sprintf("other-%d", i)
	}

	cache.locks[blocked].Lock()
	done := make(chan struct{})
	result := make(chan error, 1)
	defer func() {
		cache.locks[blocked].Unlock()
		<-done
	}()

	go func() {
		defer close(done)
		_, err := cache.SetIfAbsent(context.Background(), key, "value", time.Minute)
		result <- err
	}()

	select {
	case err := <-result:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("a blocked stripe stalled a write on another stripe")
	}
}

func TestOverlappingSetManyDoesNotDeadlock(t *testing.T) {
	cache := newTestLocalCache(t)
	const writers = 32
	results := make(chan error, writers)
	var wg sync.WaitGroup
	for i := range writers {
		wg.Go(func() {
			results <- cache.SetMany(context.Background(), map[string]string{
				"first":  fmt.Sprint(i),
				"second": fmt.Sprint(i),
			}, time.Minute)
		})
	}

	for range writers {
		select {
		case err := <-results:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatal("overlapping batches deadlocked")
		}
	}
	wg.Wait()
}
