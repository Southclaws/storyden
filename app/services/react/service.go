// Package react allows adding/removing reactions on posts.
package react

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/el-mike/restrict"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
)

type Service interface {
	Add(ctx context.Context, accountID account.AccountID, postID post.PostID, emoji string) (*react.React, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac *restrict.AccessManager

	post_repo  post.Repository
	react_repo react.Repository
}

func New(
	l *zap.Logger,
	rbac *restrict.AccessManager,

	post_repo post.Repository,
	react_repo react.Repository,
) Service {
	return &service{
		l:          l.With(zap.String("service", "post")),
		rbac:       rbac,
		post_repo:  post_repo,
		react_repo: react_repo,
	}
}

func (s *service) Add(ctx context.Context, accountID account.AccountID, postID post.PostID, emoji string) (*react.React, error) {
	r, err := s.react_repo.Add(ctx, accountID, xid.ID(postID), emoji)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return r, nil
}
