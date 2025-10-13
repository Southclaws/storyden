package thread_mark

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/dboslee/lru"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
)

var ErrInvalidThreadMark = fault.New("invalid thread mark: thread mark did not point to a valid thread ID", ftag.With(ftag.NotFound))

// from xid
const xidEncodedLength = 20

type Service interface {
	Lookup(ctx context.Context, threadmark string) (post.ID, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	cache       *lru.SyncCache[string, xid.ID]
	thread_repo *thread_querier.Querier
}

func New(
	thread_repo *thread_querier.Querier,
) Service {
	return &service{
		cache:       lru.NewSync[string, xid.ID](lru.WithCapacity(1000)),
		thread_repo: thread_repo,
	}
}

func (s *service) Lookup(ctx context.Context, threadmark string) (post.ID, error) {
	// input is too short to be anything useful
	if len(threadmark) < xidEncodedLength {
		return post.ID(xid.NilID()), ErrInvalidThreadMark
	}

	if cv, ok := s.cache.Get(threadmark); ok {
		return post.ID(cv), nil
	}

	// the input is in the format "<xid>-<thread-slug>"
	if id, err := xid.FromString(threadmark[:xidEncodedLength]); err == nil {
		return post.ID(id), nil
	}

	// doesn't currently support any other clever thread mark lookups.
	//
	// potential future support if the desire exists:
	// - lookup by only the slug
	// - slug normalisation, like Wordpress
	return post.ID(xid.NilID()), ErrInvalidThreadMark
}
