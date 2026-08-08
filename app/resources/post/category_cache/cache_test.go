package category_cache

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidateManyBatchesDeduplicatedKeysAtOneInstant(t *testing.T) {
	store := &batchRecordingStore{}
	now := time.Date(2026, time.August, 8, 12, 0, 0, 123, time.UTC)
	clockCalls := 0
	cache := &Cache{store: store, clock: func() time.Time {
		clockCalls++
		return now.Add(time.Duration(clockCalls) * time.Nanosecond)
	}}

	require.NoError(t, cache.InvalidateMany(context.Background(), "first", "second", "first", ""))

	wantValue := now.Add(time.Nanosecond).Format(storeTimeFmt)
	assert.Equal(t, map[string]string{
		cachePrefix + "first":  wantValue,
		cachePrefix + "second": wantValue,
	}, store.values)
	assert.Equal(t, 1, store.batchCalls)
	assert.Equal(t, 1, clockCalls)
}

func BenchmarkInvalidateMany(b *testing.B) {
	for _, count := range []int{1, 10, 100, 1000} {
		b.Run(fmt.Sprintf("keys-%d", count), func(b *testing.B) {
			slugs := make([]string, count)
			for i := range slugs {
				slugs[i] = fmt.Sprintf("category-%d", i)
			}
			cache := &Cache{store: &batchRecordingStore{}, clock: time.Now}
			ctx := context.Background()

			for b.Loop() {
				if err := cache.InvalidateMany(ctx, slugs...); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type batchRecordingStore struct {
	values     map[string]string
	batchCalls int
}

func (s *batchRecordingStore) SetMany(_ context.Context, values map[string]string, _ time.Duration) error {
	s.batchCalls++
	s.values = values
	return nil
}

func (*batchRecordingStore) Get(context.Context, string) (string, error) {
	return "", errors.New("not found")
}
func (*batchRecordingStore) Set(context.Context, string, string, time.Duration) error { return nil }
func (*batchRecordingStore) Delete(context.Context, string) error                     { return nil }
func (*batchRecordingStore) HIncrBy(context.Context, string, string, int64) (int, error) {
	return 0, nil
}
func (*batchRecordingStore) HGetAll(context.Context, string) (map[string]string, error) {
	return nil, nil
}
func (*batchRecordingStore) HDel(context.Context, string, string) error          { return nil }
func (*batchRecordingStore) Expire(context.Context, string, time.Duration) error { return nil }
