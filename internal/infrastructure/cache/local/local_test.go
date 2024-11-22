package local_test

import (
	"context"
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

	t.Run("concurrency", func(t *testing.T) {
		r := require.New(t)
		ctx := context.Background()

		c, err := local.New()
		r.NoError(err)

		go func() {
			for i := 0; i < 10000; i++ {
				_, err := c.HIncrBy(ctx, "key", "field", 1)
				r.NoError(err)
			}
		}()

		go func() {
			for i := 0; i < 10000; i++ {
				_, err := c.HGetAll(ctx, "key")
				r.NoError(err)
			}
		}()
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
