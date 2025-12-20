// Package reply provides APIs for managing posts within a thread.
package reply

import (
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply_querier"
	"github.com/Southclaws/storyden/app/resources/post/reply_writer"
	"github.com/Southclaws/storyden/app/resources/post/thread_cache"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/moderation/content_policy"
	"github.com/Southclaws/storyden/app/services/reply/reply_notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Partial struct {
	Content opt.Optional[datagraph.Content]
	ReplyTo opt.Optional[post.ID]
	Meta    opt.Optional[map[string]any]
}

func (p Partial) Opts() (opts []reply_writer.Option) {
	p.Content.Call(func(v datagraph.Content) { opts = append(opts, reply_writer.WithContent(v)) })
	p.ReplyTo.Call(func(v post.ID) { opts = append(opts, reply_writer.WithReplyTo(v)) })
	p.Meta.Call(func(v map[string]any) { opts = append(opts, reply_writer.WithMeta(v)) })
	return
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		reply_notify.Build(),
	)
}

type Mutator struct {
	accountQuery *account_querier.Querier
	replyQuerier *reply_querier.Querier
	replyWriter  *reply_writer.Writer
	fetcher      *fetcher.Fetcher
	bus          *pubsub.Bus
	cpm          *content_policy.Manager
	cache        *thread_cache.Cache
}

func New(
	accountQuery *account_querier.Querier,
	replyQuerier *reply_querier.Querier,
	replyWriter *reply_writer.Writer,
	fetcher *fetcher.Fetcher,
	bus *pubsub.Bus,
	cpm *content_policy.Manager,
	cache *thread_cache.Cache,
) *Mutator {
	return &Mutator{
		accountQuery: accountQuery,
		replyQuerier: replyQuerier,
		replyWriter:  replyWriter,
		fetcher:      fetcher,
		bus:          bus,
		cpm:          cpm,
		cache:        cache,
	}
}
