package thread

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_category "github.com/Southclaws/storyden/internal/ent/category"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

// Note: The resources thread and post both map to the same underlying database
// schema ent. The point of the resources being separate is to provide
// separate intuitive APIs that abstract away the detail that a `post` item in
// the database and a `thread` item use the same underlying table.

type Repository interface {
	// Create a new thread. A thread is just a "post" in the underlying data
	// ent. But a thread is marked as "first" and has a title, catgegory and
	// tags, and no parent post.
	Create(
		ctx context.Context,
		title string,
		body string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		tags []string,
		opts ...Option,
	) (*Thread, error)

	Update(ctx context.Context, id post.ID, opts ...Option) (*Thread, error)

	List(
		ctx context.Context,
		before time.Time,
		max int,
		opts ...Query,
	) ([]*Thread, error)

	// GetPostCounts(ctx context.Context) (map[string]int, error)

	Get(ctx context.Context, threadID post.ID) (*Thread, error)

	Delete(ctx context.Context, id post.ID) error
}

type Option func(*ent.PostMutation)

func WithID(id post.ID) Option {
	return func(m *ent.PostMutation) {
		m.SetID(xid.ID(id))
	}
}

func WithTitle(v string) Option {
	return func(pm *ent.PostMutation) {
		pm.SetTitle(v)
	}
}

func WithSummary(v string) Option {
	return func(pm *ent.PostMutation) {
		pm.SetShort(v)
	}
}

func WithBody(v string) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(v)
	}
}

func WithTags(v []xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddTagIDs(v...)
	}
}

func WithCategory(v xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetCategoryID(v)
	}
}

func WithStatus(v post.Status) Option {
	return func(pm *ent.PostMutation) {
		pm.SetStatus(v.ToEnt())
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.PostMutation) {
		m.SetMetadata(meta)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.PostMutation) {
		m.AddAssetIDs(dt.Map(a, func(id asset.AssetID) string { return string(id) })...)
	}
}

func WithLinks(ids ...xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddLinkIDs(ids...)
	}
}

type Query func(q *ent.PostQuery)

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

func HasStatus(status post.Status) Query {
	return func(q *ent.PostQuery) {
		q.Where(ent_post.StatusEQ(status.ToEnt()))
	}
}
