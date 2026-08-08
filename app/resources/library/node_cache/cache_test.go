package node_cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/mark"
)

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

func TestInvalidateRefreshesCurrentAndDeletesRetiredAliases(t *testing.T) {
	ctx := context.Background()
	store := newRecordingStore()
	now := time.Date(2026, time.August, 8, 12, 0, 0, 123, time.UTC)
	clockCalls := 0
	cache := &Cache{store: store, clock: func() time.Time {
		clockCalls++
		return now.Add(time.Duration(clockCalls) * time.Nanosecond)
	}}
	id := xid.New()
	store.values[cachePrefix+"old-slug"] = "old-validator"

	require.NoError(t, cache.Invalidate(ctx, id, "new-slug", "old-slug", "old-slug", "new-slug"))

	wantValue := now.Add(time.Nanosecond).Format(storeTimeFmt)
	assert.Equal(t, map[string]string{
		cachePrefix + id.String(): wantValue,
		cachePrefix + "new-slug":  wantValue,
	}, store.values)
	assert.Equal(t, 1, clockCalls, "all aliases must share one invalidation instant")
	assert.Equal(t, 2, store.setCalls)
	assert.Equal(t, 1, store.batchCalls)
}

func TestInvalidateReturnsBatchStoreFailure(t *testing.T) {
	ctx := context.Background()
	store := newRecordingStore()
	store.setErr = errors.New("cache unavailable")
	cache := &Cache{store: store, clock: time.Now}
	id := xid.New()

	err := cache.Invalidate(ctx, id, "new-slug", "old-slug")

	require.ErrorIs(t, err, store.setErr)
	assert.Empty(t, store.values)
	assert.Equal(t, 0, store.setCalls)
	assert.Equal(t, 1, store.batchCalls)
}

func TestDeleteNodeDeletesIDAndSlugAliases(t *testing.T) {
	ctx := context.Background()
	store := newRecordingStore()
	cache := &Cache{store: store, clock: time.Now}
	id := xid.New()
	store.values[cachePrefix+id.String()] = "id"
	store.values[cachePrefix+"node-slug"] = "slug"

	require.NoError(t, cache.deleteNode(ctx, id, "node-slug"))

	assert.Empty(t, store.values)
}

type recordingStore struct {
	values     map[string]string
	setCalls   int
	batchCalls int
	setErr     error
}

func newRecordingStore() *recordingStore {
	return &recordingStore{values: map[string]string{}}
}

func (s *recordingStore) Get(_ context.Context, key string) (string, error) {
	value, ok := s.values[key]
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

func (s *recordingStore) Set(_ context.Context, key string, value string, _ time.Duration) error {
	s.setCalls++
	if s.setErr != nil {
		return s.setErr
	}
	s.values[key] = value
	return nil
}

func (s *recordingStore) SetMany(_ context.Context, values map[string]string, _ time.Duration) error {
	s.batchCalls++
	if s.setErr != nil {
		return s.setErr
	}

	for key, value := range values {
		s.setCalls++
		s.values[key] = value
	}

	return nil
}

func (s *recordingStore) Delete(_ context.Context, key string) error {
	delete(s.values, key)
	return nil
}

func (s *recordingStore) HIncrBy(context.Context, string, string, int64) (int, error) {
	panic("not implemented")
}

func (s *recordingStore) HGetAll(context.Context, string) (map[string]string, error) {
	panic("not implemented")
}

func (s *recordingStore) HDel(context.Context, string, string) error {
	panic("not implemented")
}

func (s *recordingStore) Expire(context.Context, string, time.Duration) error {
	panic("not implemented")
}

func BenchmarkCanonicalKey(b *testing.B) {
	id := xid.New()

	for name, key := range map[string]mark.Queryable{
		"slug":        mark.NewQueryKey("node-slug"),
		"id":          mark.NewQueryKeyID(id),
		"id-and-slug": mark.NewQueryKey(id.String() + "-node-slug"),
	} {
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				_ = CanonicalKey(key)
			}
		})
	}
}

func BenchmarkInvalidate(b *testing.B) {
	cache := &Cache{store: noOpStore{}, clock: time.Now}
	id := xid.New()
	ctx := context.Background()

	for b.Loop() {
		if err := cache.Invalidate(ctx, id, "new-slug", "old-slug"); err != nil {
			b.Fatal(err)
		}
	}
}

type noOpStore struct{}

func (noOpStore) Get(context.Context, string) (string, error)              { return "", nil }
func (noOpStore) Set(context.Context, string, string, time.Duration) error { return nil }
func (noOpStore) SetMany(context.Context, map[string]string, time.Duration) error {
	return nil
}
func (noOpStore) Delete(context.Context, string) error                        { return nil }
func (noOpStore) HIncrBy(context.Context, string, string, int64) (int, error) { return 0, nil }
func (noOpStore) HGetAll(context.Context, string) (map[string]string, error)  { return nil, nil }
func (noOpStore) HDel(context.Context, string, string) error                  { return nil }
func (noOpStore) Expire(context.Context, string, time.Duration) error         { return nil }
