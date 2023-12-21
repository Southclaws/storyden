// Package thread provides APIs for working with threads which are sequences of
// posts. Threads can be created with one post, listed, searched and updated.
package thread

import (
	"context"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/hydrator"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type Service interface {
	// Create a new thread in the specified category.
	Create(
		ctx context.Context,
		title string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		status post.Status,
		tags []string,
		meta map[string]any,
		partial Partial,
	) (*thread.Thread, error)

	Update(ctx context.Context, threadID post.ID, partial Partial) (*thread.Thread, error)

	Delete(ctx context.Context, id post.ID) error

	// ListAll returns all threads.
	ListAll(
		ctx context.Context,
		before time.Time,
		max int,
		query Params,
	) ([]*thread.Thread, error)

	// Get one thread and the posts within it.
	Get(
		ctx context.Context,
		threadID post.ID,
	) (*thread.Thread, error)
}

type Partial struct {
	Title    opt.Optional[string]
	Body     opt.Optional[string]
	Tags     opt.Optional[[]xid.ID]
	Category opt.Optional[xid.ID]
	Status   opt.Optional[post.Status]
	URL      opt.Optional[string]
	Meta     opt.Optional[map[string]any]
}

func (p Partial) Opts() (opts []thread.Option) {
	p.Title.Call(func(v string) { opts = append(opts, thread.WithTitle(v)) })
	p.Body.Call(func(v string) { opts = append(opts, thread.WithBody(v)) })
	p.Tags.Call(func(v []xid.ID) { opts = append(opts, thread.WithTags(v)) })
	p.Category.Call(func(v xid.ID) { opts = append(opts, thread.WithCategory(xid.ID(v))) })
	p.Status.Call(func(v post.Status) { opts = append(opts, thread.WithStatus(v)) })
	p.Meta.Call(func(v map[string]any) { opts = append(opts, thread.WithMeta(v)) })
	return
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
	thread_repo  thread.Repository
	hydrator     hydrator.Service
	semdex       semdex.Service
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	thread_repo thread.Repository,
	hydrator hydrator.Service,
	semdex semdex.Service,
) Service {
	return &service{
		l:            l.With(zap.String("service", "thread")),
		rbac:         rbac,
		account_repo: account_repo,
		thread_repo:  thread_repo,
		hydrator:     hydrator,
		semdex:       semdex,
	}
}
