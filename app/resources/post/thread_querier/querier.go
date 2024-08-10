package thread_querier

import (
	"context"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_category "github.com/Southclaws/storyden/internal/ent/category"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

type Result struct {
	PageSize    int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Threads     []*thread.Thread
}

type Querier interface {
	// List is used for listing threads or filtering with basic queries. It's
	// not sufficient for full scale text-based or semantic search however.
	List(ctx context.Context,
		page int,
		size int,
		opts ...Query,
	) (*Result, error)

	Get(ctx context.Context, threadID post.ID) (*thread.Thread, error)
}

type Query func(q *ent.PostQuery)

func HasKeyword(s string) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.Or(
			ent_post.TitleContainsFold(s),
			ent_post.SlugContainsFold(s),
			ent_post.BodyContainsFold(s),
		))
	}
}

func HasCreatedDateBefore(t time.Time) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.CreatedAtLT(t))
	}
}

func HasUpdatedDateBefore(t time.Time) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.UpdatedAtLT(t))
	}
}

func HasAuthor(id account.AccountID) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.HasAuthorWith(ent_account.ID(xid.ID(id))))
	}
}

func HasTags(ids []xid.ID) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.HasTagsWith(ent_tag.IDIn(ids...)))
	}
}

func HasCategories(ids []string) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.HasCategoryWith(ent_category.SlugIn(ids...)))
	}
}

func HasStatus(status visibility.Visibility) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.VisibilityEQ(ent_post.Visibility(status.String())))
	}
}

func HasNotBeenDeleted() Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.DeletedAtIsNil())
	}
}
