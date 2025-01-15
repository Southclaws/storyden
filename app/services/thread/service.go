// Package thread provides APIs for working with threads which are sequences of
// posts. Threads can be created with one post, listed, searched and updated.
package thread

import (
	"context"
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/tag/tag_writer"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/mention/mentioner"
	"github.com/Southclaws/storyden/app/services/moderation/content_policy"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/thread/thread_semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Service interface {
	// Create a new thread in the specified category.
	Create(
		ctx context.Context,
		title string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		status visibility.Visibility,
		meta map[string]any,
		partial Partial,
	) (*thread.Thread, error)

	Update(ctx context.Context, threadID post.ID, partial Partial) (*thread.Thread, error)

	Delete(ctx context.Context, id post.ID) error

	List(ctx context.Context,
		page int,
		size int,
		opts Params,
	) (*thread.Result, error)

	// Get one thread and the posts within it.
	Get(
		ctx context.Context,
		threadID post.ID,
		pageParams pagination.Parameters,
	) (*thread.Thread, error)
}

type Partial struct {
	Title      opt.Optional[string]
	Content    opt.Optional[datagraph.Content]
	Category   opt.Optional[xid.ID]
	Tags       opt.Optional[tag_ref.Names]
	Visibility opt.Optional[visibility.Visibility]
	URL        opt.Optional[url.URL]
	Meta       opt.Optional[map[string]any]
}

func (p Partial) Opts() (opts []thread.Option) {
	p.Title.Call(func(v string) { opts = append(opts, thread.WithTitle(v)) })
	p.Content.Call(func(v datagraph.Content) { opts = append(opts, thread.WithContent(v)) })
	p.Category.Call(func(v xid.ID) { opts = append(opts, thread.WithCategory(xid.ID(v))) })
	p.Visibility.Call(func(v visibility.Visibility) { opts = append(opts, thread.WithVisibility(v)) })
	p.Meta.Call(func(v map[string]any) { opts = append(opts, thread.WithMeta(v)) })
	return
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		thread_semdex.Build(),
	)
}

type service struct {
	ins spanner.Instrumentation
	l   *zap.Logger

	accountQuery *account_querier.Querier
	thread_repo  thread.Repository
	tagWriter    *tag_writer.Writer
	fetcher      *fetcher.Fetcher
	recommender  semdex.Recommender
	indexQueue   pubsub.Topic[mq.IndexThread]
	deleteQueue  pubsub.Topic[mq.DeleteThread]
	mentioner    *mentioner.Mentioner
	cpm          *content_policy.Manager
}

func New(
	ins spanner.Builder,
	l *zap.Logger,

	accountQuery *account_querier.Querier,
	thread_repo thread.Repository,
	tagWriter *tag_writer.Writer,
	fetcher *fetcher.Fetcher,
	recommender semdex.Recommender,
	indexQueue pubsub.Topic[mq.IndexThread],
	deleteQueue pubsub.Topic[mq.DeleteThread],
	mentioner *mentioner.Mentioner,
	cpm *content_policy.Manager,
) Service {
	return &service{
		ins: ins.Build(),
		l:   l.With(zap.String("service", "thread")),

		accountQuery: accountQuery,
		thread_repo:  thread_repo,
		tagWriter:    tagWriter,
		fetcher:      fetcher,
		recommender:  recommender,
		indexQueue:   indexQueue,
		deleteQueue:  deleteQueue,
		mentioner:    mentioner,
		cpm:          cpm,
	}
}
