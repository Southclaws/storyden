// Package thread provides APIs for working with threads which are sequences of
// posts. Threads can be created with one post, listed, searched and updated.
package thread

import (
	"context"
	"net/url"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/post/thread_cache"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/app/resources/post/thread_writer"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/tag/tag_writer"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/mention/mentioner"
	"github.com/Southclaws/storyden/app/services/moderation"
	"github.com/Southclaws/storyden/app/services/report/system_report"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Service interface {
	// Create a new thread with optional category.
	Create(
		ctx context.Context,
		title string,
		authorID account.AccountID,
		meta map[string]any,
		partial Partial,
	) (*thread.Thread, error)

	Update(ctx context.Context, threadID post.ID, partial Partial) (*thread.Thread, error)

	Delete(ctx context.Context, id post.ID) error

	List(ctx context.Context,
		page int,
		size int,
		opts Params,
	) (*thread_querier.Result, error)

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
	Pinned     opt.Optional[int]
}

func (p Partial) Opts() (opts []thread_writer.Option) {
	p.Title.Call(func(v string) { opts = append(opts, thread_writer.WithTitle(v)) })
	p.Content.Call(func(v datagraph.Content) { opts = append(opts, thread_writer.WithContent(v)) })
	p.Category.Call(func(v xid.ID) { opts = append(opts, thread_writer.WithCategory(xid.ID(v))) })
	p.Visibility.Call(func(v visibility.Visibility) { opts = append(opts, thread_writer.WithVisibility(v)) })
	p.Meta.Call(func(v map[string]any) { opts = append(opts, thread_writer.WithMeta(v)) })
	p.Pinned.Call(func(v int) { opts = append(opts, thread_writer.WithPinned(v)) })
	return
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
	)
}

type service struct {
	ins spanner.Instrumentation

	accountQuery   *account_querier.Querier
	threadQuerier  *thread_querier.Querier
	threadWriter   *thread_writer.Writer
	tagWriter      *tag_writer.Writer
	fetcher        *fetcher.Fetcher
	recommender    semdex.Recommender
	bus            *pubsub.Bus
	mentioner      *mentioner.Mentioner
	cpm            *moderation.Manager
	cache          *thread_cache.Cache
	systemReporter *system_report.Manager
}

func New(
	ins spanner.Builder,

	accountQuery *account_querier.Querier,
	threadQuerier *thread_querier.Querier,
	threadWriter *thread_writer.Writer,
	tagWriter *tag_writer.Writer,
	fetcher *fetcher.Fetcher,
	recommender semdex.Recommender,
	bus *pubsub.Bus,
	mentioner *mentioner.Mentioner,
	cpm *moderation.Manager,
	cache *thread_cache.Cache,
	systemReporter *system_report.Manager,
) Service {
	return &service{
		ins: ins.Build(),

		accountQuery:   accountQuery,
		threadQuerier:  threadQuerier,
		threadWriter:   threadWriter,
		tagWriter:      tagWriter,
		fetcher:        fetcher,
		recommender:    recommender,
		bus:            bus,
		mentioner:      mentioner,
		cpm:            cpm,
		cache:          cache,
		systemReporter: systemReporter,
	}
}
