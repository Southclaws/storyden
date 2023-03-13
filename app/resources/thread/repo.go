package thread

import (
	"context"
	"time"

	"github.com/rs/xid"

	account_resource "github.com/Southclaws/storyden/app/resources/account"
	category_resource "github.com/Southclaws/storyden/app/resources/category"
	post_resource "github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/category"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/tag"
)

// Note: The resources thread and post both map to the same underlying database
// schema ent. The point of the resources being separate is to provide
// separate intuitive APIs that abstract away the detail that a `post` item in
// the database and a `thread` item use the same underlying table.

type option func(*Thread)

type Repository interface {
	// Create a new thread. A thread is just a "post" in the underlying data
	// ent. But a thread is marked as "first" and has a title, catgegory and
	// tags, and no parent post.
	Create(
		ctx context.Context,
		title string,
		body string,
		authorID account_resource.AccountID,
		categoryID category_resource.CategoryID,
		tags []string,
		opts ...option,
	) (*Thread, error)

	List(
		ctx context.Context,
		before time.Time,
		max int,
		opts ...Query,
	) ([]*Thread, error)

	// GetPostCounts(ctx context.Context) (map[string]int, error)

	Get(ctx context.Context, threadID post_resource.PostID) (*Thread, error)

	// Update(ctx context.Context, userID user.UserID, id string, title, category *string, pinned *bool) (*post.Post, error)

	// Delete(ctx context.Context, id, authorID user.UserID) (int, error)
}

func WithID(id post_resource.PostID) option {
	return func(c *Thread) {
		c.ID = id
	}
}

func WithMeta(meta map[string]any) option {
	return func(t *Thread) {
		t.Meta = meta
	}
}

type Query func(q *ent.PostQuery)

func WithAuthor(id account_resource.AccountID) Query {
	return func(q *ent.PostQuery) {
		q.Where(post.HasAuthorWith(account.ID(xid.ID(id))))
	}
}

func WithTags(ids []xid.ID) Query {
	return func(q *ent.PostQuery) {
		q.Where(post.HasTagsWith(tag.IDIn(ids...)))
	}
}

func WithCategories(ids []string) Query {
	return func(q *ent.PostQuery) {
		q.Where(post.HasCategoryWith(category.NameIn(ids...)))
	}
}
