package post

import (
	"context"

	"4d63.com/optional"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type Option func(*ent.PostMutation)

type Repository interface {
	Create(
		ctx context.Context,
		body string,
		authorID account.AccountID,
		parentID PostID,
		replyToID optional.Optional[PostID],
		meta map[string]any,
		opts ...Option,
	) (*Post, error)

	// EditPost(ctx context.Context, authorID, postID string, title *string, body *string) (*Post, error)
	// DeletePost(ctx context.Context, authorID, postID string, force bool) (*Post, error)
}

func WithID(id PostID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetID(xid.ID(id))
	}
}
