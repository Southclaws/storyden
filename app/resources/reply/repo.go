package reply

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type Option func(*ent.PostMutation)

type Repository interface {
	Create(
		ctx context.Context,
		body string,
		authorID account.AccountID,
		parentID post.ID,
		replyToID opt.Optional[post.ID],
		meta map[string]any,
		opts ...Option,
	) (*Reply, error)

	Get(ctx context.Context, id post.ID) (*Reply, error)

	Update(ctx context.Context, id post.ID, opts ...Option) (*Reply, error)
	// EditPost(ctx context.Context, authorID, postID string, title *string, body *string) (*Post, error)
	Delete(ctx context.Context, id post.ID) error
}

func WithID(id post.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetID(xid.ID(id))
	}
}

func WithBody(v string) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(string(v))
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.PostMutation) {
		m.SetMetadata(meta)
	}
}

func WithAssets(ids ...asset.AssetID) Option {
	return func(m *ent.PostMutation) {
		m.AddAssetIDs(dt.Map(ids, func(id asset.AssetID) string { return string(id) })...)
	}
}

func WithLinks(ids ...xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddLinkIDs(ids...)
	}
}
