// Package reply provides APIs for managing posts within a thread.
package reply

import (
	"context"

	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/notification/notify"
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

	accountQuery account_querier.Querier
	post_repo    reply.Repository
	fetcher      *fetcher.Fetcher
	indexQueue   pubsub.Topic[mq.IndexPost]
	notifier     *notify.Notifier
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	accountQuery account_querier.Querier,
	post_repo reply.Repository,
	fetcher *fetcher.Fetcher,
	indexQueue pubsub.Topic[mq.IndexPost],
	notifier *notify.Notifier,
) Service {
	return &service{
		l:            l.With(zap.String("service", "reply")),
		rbac:         rbac,
		accountQuery: accountQuery,
		post_repo:    post_repo,
		fetcher:      fetcher,
		indexQueue:   indexQueue,
		notifier:     notifier,
	}
}
