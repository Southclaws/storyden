package thread_writer

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Repository interface {
	Create(
		ctx context.Context,
		title string,
		authorID account.AccountID,
		categoryID category.CategoryID,
		tags []string,
		opts ...Option,
	) (*thread.Thread, error)

	Update(ctx context.Context, id post.ID, opts ...Option) (*thread.Thread, error)

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

func WithContent(v content.Rich) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(v.HTML())
		pm.SetShort(v.Short())
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

func WithVisibility(v visibility.Visibility) Option {
	return func(pm *ent.PostMutation) {
		pm.SetVisibility(ent_post.Visibility(v.String()))
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.PostMutation) {
		m.SetMetadata(meta)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.PostMutation) {
		m.AddAssetIDs(a...)
	}
}

func WithLink(id xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetLinkID(id)
	}
}

func WithContentLinks(ids ...xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddContentLinkIDs(ids...)
	}
}
