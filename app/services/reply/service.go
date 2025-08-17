// Package reply provides APIs for managing posts within a thread.
package reply

import (
	"context"

	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/moderation/content_policy"
	"github.com/Southclaws/storyden/app/services/reply/reply_notify"
	"github.com/Southclaws/storyden/app/services/reply/reply_semdex"
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
	Content opt.Optional[datagraph.Content]
	ReplyTo opt.Optional[post.ID]
	Meta    opt.Optional[map[string]any]
}

func (p Partial) Opts() (opts []reply.Option) {
	p.Content.Call(func(v datagraph.Content) { opts = append(opts, reply.WithContent(v)) })
	p.ReplyTo.Call(func(v post.ID) { opts = append(opts, reply.WithReplyTo(v)) })
	p.Meta.Call(func(v map[string]any) { opts = append(opts, reply.WithMeta(v)) })
	return
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		reply_semdex.Build(),
		reply_notify.Build(),
	)
}

type service struct {
	accountQuery *account_querier.Querier
	post_repo    reply.Repository
	fetcher      *fetcher.Fetcher
	bus          *pubsub.Bus
	cpm          *content_policy.Manager
}

func New(
	accountQuery *account_querier.Querier,
	post_repo reply.Repository,
	fetcher *fetcher.Fetcher,
	bus *pubsub.Bus,
	cpm *content_policy.Manager,
) Service {
	return &service{
		accountQuery: accountQuery,
		post_repo:    post_repo,
		fetcher:      fetcher,
		bus:          bus,
		cpm:          cpm,
	}
}
