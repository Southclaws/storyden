package thread_mark

import (
	"context"
	"log/slog"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

var ErrInvalidThreadMark = fault.New("invalid thread mark: thread mark did not point to a valid thread ID", ftag.With(ftag.NotFound))

const xidEncodedLength = 20

type Service interface {
	Lookup(ctx context.Context, threadmark string) (post.ID, error)
}

func Build() fx.Option { return fx.Provide(New) }

type service struct {
	cache  cache.Store
	logger *slog.Logger
}

func New(store cache.Store, logger *slog.Logger) Service {
	return &service{cache: store, logger: logger}
}

func (s *service) Lookup(ctx context.Context, threadmark string) (post.ID, error) {
	if len(threadmark) < xidEncodedLength {
		return post.ID(xid.NilID()), ErrInvalidThreadMark
	}
	prefix := threadmark[:xidEncodedLength]
	key := "thread:mark:" + prefix
	if stored, err := s.cache.Get(ctx, key); err == nil && stored == prefix {
		if id, err := xid.FromString(stored); err == nil {
			return post.ID(id), nil
		}
	}
	id, err := xid.FromString(prefix)
	if err != nil {
		return post.ID(xid.NilID()), ErrInvalidThreadMark
	}
	if err := s.cache.Set(ctx, key, id.String(), time.Hour); err != nil {
		s.logger.WarnContext(ctx, "failed to cache thread mark", slog.String("error", err.Error()))
	}
	return post.ID(id), nil
}
