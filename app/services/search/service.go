package search

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/services/search/simplesearch"
)

type Service interface {
	Search(ctx context.Context, q Query) ([]*post.Post, error)
}

type Query struct {
	Body   opt.Optional[string]
	Author opt.Optional[string]
	Kinds  opt.Optional[[]post_search.Kind]
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Provide(simplesearch.NewParallelSearcher, NewSearcher),
	)
}

type service struct {
	l *zap.Logger

	post_search_repo post_search.Repository
}

func New(
	l *zap.Logger,

	accountQuery *account_querier.Querier,
	post_search_repo post_search.Repository,
) Service {
	return &service{
		l: l.With(zap.String("service", "search")),

		post_search_repo: post_search_repo,
	}
}

func (s *service) Search(ctx context.Context, q Query) ([]*post.Post, error) {
	filters := []post_search.Filter{}

	q.Body.Call(func(v string) {
		filters = append(filters, post_search.WithKeywords(v))
	})
	q.Author.Call(func(v string) {
		filters = append(filters, post_search.WithAuthorHandle(v))
	})
	q.Kinds.Call(func(v []post_search.Kind) {
		filters = append(filters, post_search.WithKinds(v...))
	})

	posts, err := s.post_search_repo.Search(ctx, filters...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return posts, nil
}
