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
)

type Service interface {
	// Create a new thread in the specified category.
	Create(
		ctx context.Context,
		title string,
		body string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		status thread.Status,
		tags []string,
		meta map[string]any,
		opts ...thread.Option,
	) (*thread.Thread, error)

	Update(ctx context.Context, threadID post.PostID, partial Partial) (*thread.Thread, error)

	Delete(ctx context.Context, id post.PostID) error

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
		threadID post.PostID,
	) (*thread.Thread, error)
}

type Partial struct {
	Title    opt.Optional[string]
	Body     opt.Optional[string]
	Tags     opt.Optional[[]xid.ID]
	Category opt.Optional[xid.ID]
	Status   opt.Optional[thread.Status]
	Meta     opt.Optional[map[string]any]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo account.Repository
	thread_repo  thread.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	thread_repo thread.Repository,
) Service {
	return &service{
		l:            l.With(zap.String("service", "thread")),
		rbac:         rbac,
		account_repo: account_repo,
		thread_repo:  thread_repo,
	}
}
