package local_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/infrastructure/cache/local"
)

// NOTE: Tiny sleeps between set/get because ristretto is eventually consistent.

func TestLocalCache(t *testing.T) {
	t.Run("get_set", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		key := "key"
		value := "value"

		err = c.Set(ctx, key, value, time.Minute)
		r.NoError(err)

		// SEE NOTE
		time.Sleep(time.Millisecond)

		v, err := c.Get(ctx, key)
		r.NoError(err)

		a.Equal(value, v)
	})

	t.Run("hincr", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		key := "key"
		field := "f"

		i, err := c.HIncrBy(ctx, key, field, 1)
		r.NoError(err)
		a.Equal(1, i)

		// SEE NOTE
		time.Sleep(time.Millisecond)

		i, err = c.HIncrBy(ctx, key, field, 1)
		r.NoError(err)
		a.Equal(2, i)

		m, err := c.HGetAll(ctx, key)
		r.NoError(err)
		a.Equal(map[string]string{field: "2"}, m)
	})

	t.Run("set_many", func(t *testing.T) {
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		values := map[string]string{"first": "one", "second": "two", "third": "three"}
		r.NoError(c.SetMany(ctx, values, time.Minute))

		for key, want := range values {
			got, err := c.Get(ctx, key)
			r.NoError(err)
			r.Equal(want, got)
		}
	})

	t.Run("set_if_absent_concurrent_claims", func(t *testing.T) {
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		key := "key"
		const count = 32
		inserted := make([]bool, count)
		errors := make([]error, count)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := range count {
			wg.Go(func() {
				<-start
				inserted[i], errors[i] = c.SetIfAbsent(ctx, key, fmt.Sprint(i), time.Minute)
			})
		}
		close(start)
		wg.Wait()

		winners := 0
		for i, won := range inserted {
			r.NoError(errors[i])
			if won {
				winners++
				value, err := c.Get(ctx, key)
				r.NoError(err)
				r.Equal(fmt.Sprint(i), value)
			}
		}
		r.Equal(1, winners)
	})

	t.Run("set_if_absent_expiry_is_not_extended_by_duplicates", func(t *testing.T) {
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		key := "key"
		inserted, err := c.SetIfAbsent(ctx, key, "first", 250*time.Millisecond)
		r.NoError(err)
		r.True(inserted)
		inserted, err = c.SetIfAbsent(ctx, key, "second", time.Hour)
		r.NoError(err)
		r.False(inserted)

		r.Eventually(func() bool {
			_, err := c.Get(ctx, key)
			return err != nil
		}, 3*time.Second, 10*time.Millisecond)

		inserted, err = c.SetIfAbsent(ctx, key, "third", time.Minute)
		r.NoError(err)
		r.True(inserted)
		value, err := c.Get(ctx, key)
		r.NoError(err)
		r.Equal("third", value)
	})

	t.Run("set_if_absent_existing_keys_and_delete", func(t *testing.T) {
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		key := "key"
		r.NoError(c.Set(ctx, key, "existing", time.Minute))
		inserted, err := c.SetIfAbsent(ctx, key, "replacement", time.Minute)
		r.NoError(err)
		r.False(inserted)
		value, err := c.Get(ctx, key)
		r.NoError(err)
		r.Equal("existing", value)

		r.NoError(c.Delete(ctx, key))
		inserted, err = c.SetIfAbsent(ctx, key, "replacement", time.Minute)
		r.NoError(err)
		r.True(inserted)
	})

	t.Run("set_if_absent_invalid_ttl", func(t *testing.T) {
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		for _, ttl := range []time.Duration{0, -time.Second} {
			inserted, err := c.SetIfAbsent(ctx, "key", "value", ttl)
			r.Error(err)
			r.False(inserted)
		}
	})
}

func BenchmarkLocalCache(b *testing.B) {
	// A rouch concurrency smash test to make sure we're race-free

	ctx := context.Background()

	c, err := local.New()
	if err != nil {
		panic(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.HIncrBy(ctx, "key", "field", 1)
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.HGetAll(ctx, "key")
			if err != nil {
				b.Error(err)
			}
		}
	})
}
