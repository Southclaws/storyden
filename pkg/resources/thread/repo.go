package thread

import (
	"context"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/category"
)

//go:generate mockery --inpackage --name=Repository --case=underscore

// Note: The resources thread and post both map to the same underlying database
// schema model. The point of the resources being separate is to provide
// separate intuitive APIs that abstract away the detail that a `post` item in
// the database and a `thread` item use the same underlying table.

type Repository interface {
	CreateThread(
		ctx context.Context,
		title string,
		body string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		tags []string,
	) (*Thread, error)

	// GetThreads(
	// 	ctx context.Context,
	// 	tags []string, category string, query string,
	// 	before time.Time, sort string, offset, max int,
	// 	includePosts bool,
	// 	includeDeleted bool,
	// 	includeAdmin bool,
	// ) ([]post.Post, error)

	// GetPostCounts(ctx context.Context) (map[string]int, error)

	// GetThread(ctx context.Context, slug string, max, skip int, deleted, admin bool) (Thread, error)

	// Update(ctx context.Context, userID user.UserID, id string, title, category *string, pinned *bool) (*post.Post, error)

	// Delete(ctx context.Context, id, authorID user.UserID) (int, error)
}
