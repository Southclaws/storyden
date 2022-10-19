// Package thread provides APIs for working with threads which are sequences of
// posts. Threads can be created with one post, listed, searched and updated.
package thread

import (
	"context"
	"time"

	"github.com/el-mike/restrict"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
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
		tags []string,
		meta map[string]any,
	) (*thread.Thread, error)

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

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac *restrict.AccessManager

	account_repo account.Repository
	thread_repo  thread.Repository
}

func New(
	l *zap.Logger,
	rbac *restrict.AccessManager,

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
