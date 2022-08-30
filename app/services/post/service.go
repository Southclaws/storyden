// Package post provides APIs for managing posts within a thread.
package post

import (
	"context"

	"4d63.com/optional"
	"github.com/el-mike/restrict"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
)

type Service interface {
	// Create a new thread in the specified category.
	Create(
		ctx context.Context,
		body string,
		authorID account.AccountID,
		parentID post.PostID,
		replyToID optional.Optional[post.PostID],
	) (*post.Post, error)
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac *restrict.AccessManager

	account_repo account.Repository
	post_repo    post.Repository
}

func New(
	l *zap.Logger,
	rbac *restrict.AccessManager,

	account_repo account.Repository,
	post_repo post.Repository,
) Service {
	return &service{
		l:            l.With(zap.String("service", "post")),
		rbac:         rbac,
		account_repo: account_repo,
		post_repo:    post_repo,
	}
}
