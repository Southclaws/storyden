package post

import (
	"context"

	"4d63.com/optional"

	"github.com/Southclaws/storyden/app/resources/account"
)

type Repository interface {
	Create(
		ctx context.Context,
		body string,
		authorID account.AccountID,
		parentID PostID,
		replyToID optional.Optional[PostID],
		meta map[string]any,
	) (*Post, error)

	// EditPost(ctx context.Context, authorID, postID string, title *string, body *string) (*Post, error)
	// DeletePost(ctx context.Context, authorID, postID string, force bool) (*Post, error)
}
