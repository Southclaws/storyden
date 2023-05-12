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

	Get(ctx context.Context, id PostID) (*Post, error)

	Update(ctx context.Context, id PostID, opts ...Option) (*Post, error)
	// EditPost(ctx context.Context, authorID, postID string, title *string, body *string) (*Post, error)
	Delete(ctx context.Context, id PostID) error
}

func WithID(id PostID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetID(xid.ID(id))
	}
}

func WithBody(v string) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(v)
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.PostMutation) {
		m.SetMetadata(meta)
	}
}
