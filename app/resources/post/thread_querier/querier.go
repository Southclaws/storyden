package thread_querier

import (
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_category "github.com/Southclaws/storyden/internal/ent/category"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Querier struct {
	ins         spanner.Instrumentation
	db          *ent.Client
	raw         *sqlx.DB
	roleQuerier *role_repo.Repository
}

func New(ins spanner.Builder, db *ent.Client, raw *sqlx.DB, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{
		ins:         ins.Build(),
		db:          db,
		raw:         raw,
		roleQuerier: roleQuerier,
	}
}

type Result struct {
	PageSize    int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Threads     []*thread.Thread
}

// 3 states:
// 1. Slugs filled - filter by slugs, ignore other fields.
// 2. Slugs empty, Uncategorised true - fetch uncategorised threads only.
// 3. Slugs empty, Uncategorised false - fetch all threads.
type CategoryFilter struct {
	Slugs         []string
	Uncategorised bool
}

type threadListOptions struct {
	q            *ent.PostQuery
	ignorePinned bool
}

type Query func(*threadListOptions)

func HasKeyword(s string) Query {
	return func(q *threadListOptions) {
		q.q.Where(ent_post.Or(
			ent_post.TitleContainsFold(s),
			ent_post.SlugContainsFold(s),
			ent_post.BodyContainsFold(s),
		))
	}
}

func HasCreatedDateBefore(t time.Time) Query {
	return func(q *threadListOptions) {
		q.q.Where(ent_post.CreatedAtLT(t))
	}
}

func HasUpdatedDateBefore(t time.Time) Query {
	return func(q *threadListOptions) {
		q.q.Where(ent_post.UpdatedAtLT(t))
	}
}

func HasAuthor(id account.AccountID) Query {
	return func(q *threadListOptions) {
		q.q.Where(ent_post.HasAuthorWith(ent_account.ID(xid.ID(id))))
	}
}

func HasTags(ids []xid.ID) Query {
	return func(q *threadListOptions) {
		q.q.Where(ent_post.HasTagsWith(ent_tag.IDIn(ids...)))
	}
}

func HasCategories(cf CategoryFilter) Query {
	return func(q *threadListOptions) {
		if len(cf.Slugs) > 0 {
			q.q.Where(ent_post.HasCategoryWith(ent_category.SlugIn(cf.Slugs...)))
		} else {
			if cf.Uncategorised {
				q.q.Where(ent_post.CategoryIDIsNil())
			} else {
				// No filter, fetch all threads.
			}
		}
	}
}

func HasStatus(status ...visibility.Visibility) Query {
	pv := dt.Map(status, func(v visibility.Visibility) ent_post.Visibility { return ent_post.Visibility(v.String()) })
	return func(q *threadListOptions) {
		q.q.Where(ent_post.VisibilityIn(pv...))
	}
}

func HasPublishedOrOwnInReview(accountID opt.Optional[account.AccountID], isModerator bool) Query {
	return func(q *threadListOptions) {
		publishedStatus := ent_post.Visibility(visibility.VisibilityPublished.String())
		reviewStatus := ent_post.Visibility(visibility.VisibilityReview.String())

		authorID, hasAuthor := accountID.Get()
		if !hasAuthor {
			q.q.Where(ent_post.VisibilityEQ(publishedStatus))
			return
		}

		if isModerator {
			q.q.Where(ent_post.Or(
				ent_post.VisibilityEQ(publishedStatus),
				ent_post.VisibilityEQ(reviewStatus),
			))
			return
		}

		q.q.Where(ent_post.Or(
			ent_post.VisibilityEQ(publishedStatus),
			ent_post.And(
				ent_post.VisibilityEQ(reviewStatus),
				ent_post.HasAuthorWith(ent_account.ID(xid.ID(authorID))),
			),
		))
	}
}

func HasNotBeenDeleted() Query {
	return func(q *threadListOptions) {
		q.q.Where(ent_post.DeletedAtIsNil())
	}
}

func HasNoPinnedOrdering(ignorePinned bool) Query {
	return func(q *threadListOptions) {
		q.ignorePinned = ignorePinned
	}
}
