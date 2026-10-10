package thread_mark

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"
)

type markStore struct {
	*cachetest.Store
	fail   bool
	writes int
}

func (s *markStore) Get(ctx context.Context, key string) (string, error) {
	if s.fail {
		return "", errors.New("offline")
	}
	return s.Store.Get(ctx, key)
}
func (s *markStore) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	s.writes++
	if s.fail {
		return errors.New("offline")
	}
	return s.Store.Set(ctx, key, value, ttl)
}

func TestLookup(t *testing.T) {
	t.Parallel()
	const idText = "cv1l2p1cpetc0gm4nqm0"
	id, err := xid.FromString(idText)
	require.NoError(t, err)
	for _, tc := range []struct {
		name, stored string
		offline      bool
	}{
		{"miss", "", false},
		{"hit", idText, false},
		{"corrupt", "broken", false},
		{"wrong id", "cv1l2p1cpetc0gm4nqm1", false},
		{"offline", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &markStore{Store: cachetest.New(), fail: tc.offline}
			if tc.stored != "" {
				require.NoError(t, store.Store.Set(t.Context(), "thread:mark:"+idText, tc.stored, time.Hour))
			}
			service := New(store, slog.New(slog.NewTextHandler(io.Discard, nil)))
			got, err := service.Lookup(t.Context(), idText+"-hello-world")
			require.NoError(t, err)
			require.Equal(t, post.ID(id), got)
			if tc.name == "hit" {
				require.Zero(t, store.writes)
			} else {
				require.Equal(t, 1, store.writes)
			}
			if !tc.offline {
				value, err := store.Get(t.Context(), "thread:mark:"+idText)
				require.NoError(t, err)
				require.Equal(t, idText, value)
			}
			for _, invalid := range []string{"", "short", "!!!!!!!!!!!!!!!!!!!!-slug"} {
				_, err := service.Lookup(t.Context(), invalid)
				require.ErrorIs(t, err, ErrInvalidThreadMark)
			}
		})
	}
}
