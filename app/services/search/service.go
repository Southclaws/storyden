package search

import (
	"context"

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

// TODO: Remove this entire endpoint
func (s *service) Search(ctx context.Context, q Query) ([]*post.Post, error) {
	return nil, nil
}
