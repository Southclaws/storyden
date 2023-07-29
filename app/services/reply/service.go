// Package reply provides APIs for managing posts within a thread.
package reply

import (
	"context"

	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/reply"
)

type Service interface {
	// Create a new thread in the specified category.
	Create(
		ctx context.Context,
		body string,
		authorID account.AccountID,
		parentID post.ID,
		replyToID opt.Optional[post.ID],
		meta map[string]any,
		opts ...reply.Option,
	) (*reply.Reply, error)

	Update(ctx context.Context, threadID post.ID, partial Partial) (*reply.Reply, error)

	Delete(ctx context.Context, postID post.ID) error
}

type Partial struct {
	Body opt.Optional[string]
	Meta opt.Optional[map[string]any]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
	post_repo    reply.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	post_repo reply.Repository,
) Service {
	return &service{
		l:            l.With(zap.String("service", "reply")),
		rbac:         rbac,
		account_repo: account_repo,
		post_repo:    post_repo,
	}
}
