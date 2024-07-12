package search

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type Service interface {
	Search(ctx context.Context, q Query) ([]*reply.Reply, error)
}

type Query struct {
	Body   opt.Optional[string]
	Author opt.Optional[string]
	Kinds  opt.Optional[[]post_search.Kind]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	post_search_repo post_search.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	post_search_repo post_search.Repository,
) Service {
	return &service{
		l:                l.With(zap.String("service", "search")),
		rbac:             rbac,
		post_search_repo: post_search_repo,
	}
}

func (s *service) Search(ctx context.Context, q Query) ([]*reply.Reply, error) {
	filters := []post_search.Filter{}

	q.Body.Call(func(v string) {
		filters = append(filters, post_search.WithBodyContains(v))
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
