// Package reply provides APIs for managing posts within a thread.
package reply

import (
	"context"

	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/hydrator"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Service interface {
	// Create a new thread in the specified category.
	Create(
		ctx context.Context,
		authorID account.AccountID,
		parentID post.ID,
		partial Partial,
	) (*reply.Reply, error)

	Update(ctx context.Context, threadID post.ID, partial Partial) (*reply.Reply, error)

	Delete(ctx context.Context, postID post.ID) error
}

type Partial struct {
	Content opt.Optional[content.Rich]
	ReplyTo opt.Optional[post.ID]
	Meta    opt.Optional[map[string]any]
}

func (p Partial) Opts() (opts []reply.Option) {
	p.Content.Call(func(v content.Rich) { opts = append(opts, reply.WithContent(v)) })
	p.ReplyTo.Call(func(v post.ID) { opts = append(opts, reply.WithReplyTo(v)) })
	p.Meta.Call(func(v map[string]any) { opts = append(opts, reply.WithMeta(v)) })
	return
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
	post_repo    reply.Repository
	hydrator     hydrator.Service
	indexQueue   pubsub.Topic[mq.IndexPost]
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	post_repo reply.Repository,
	hydrator hydrator.Service,
	indexQueue pubsub.Topic[mq.IndexPost],
) Service {
	return &service{
		l:            l.With(zap.String("service", "reply")),
		rbac:         rbac,
		account_repo: account_repo,
		post_repo:    post_repo,
		hydrator:     hydrator,
		indexQueue:   indexQueue,
	}
}

func (s *service) hydrate(ctx context.Context, partial Partial) (opts []reply.Option) {
	body, bodyOK := partial.Content.Get()
	if !bodyOK {
		return
	}

	return s.hydrator.HydrateReply(ctx, body, opt.NewEmpty[string]())
}
