package thread

import (
	"context"
	"time"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/post"
)

// Note: The resources thread and post both map to the same underlying database
// schema model. The point of the resources being separate is to provide
// separate intuitive APIs that abstract away the detail that a `post` item in
// the database and a `thread` item use the same underlying table.

type option func(*Thread)

type Repository interface {
	// Create a new thread. A thread is just a "post" in the underlying data
	// model. But a thread is marked as "first" and has a title, catgegory and
	// tags, and no parent post.
	Create(
		ctx context.Context,
		title string,
		body string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		tags []string,
		opts ...option,
	) (*Thread, error)

	List(
		ctx context.Context,
		before time.Time,
		max int,
	) ([]*Thread, error)

	// GetPostCounts(ctx context.Context) (map[string]int, error)

	Get(ctx context.Context, threadID post.PostID) (*Thread, error)

	// Update(ctx context.Context, userID user.UserID, id string, title, category *string, pinned *bool) (*post.Post, error)

	// Delete(ctx context.Context, id, authorID user.UserID) (int, error)
}

func WithID(id post.PostID) option {
	return func(c *Thread) {
		c.ID = id
	}
}
