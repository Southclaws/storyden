package item

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/item"
)

type (
	Option func(*ent.ItemMutation)
	Filter func(*ent.ItemQuery)
)

type Repository interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		desc string,
		opts ...Option,
	) (*datagraph.Item, error)

	Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)

	Update(ctx context.Context, slug datagraph.ItemID, opts ...Option) (*datagraph.Item, error)

	Delete(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)
}

func WithID(id datagraph.ItemID) Option {
	return func(c *ent.ItemMutation) {
		c.SetID(xid.ID(id))
	}
}

func WithName(v string) Option {
	return func(c *ent.ItemMutation) {
		c.SetName(v)
	}
}

func WithSlug(v string) Option {
	return func(c *ent.ItemMutation) {
		c.SetSlug(v)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.ItemMutation) {
		m.AddAssetIDs(a...)
	}
}

func WithAssetsRemoved(a []asset.AssetID) Option {
	return func(m *ent.ItemMutation) {
		m.RemoveAssetIDs(a...)
	}
}

func WithLinks(ids ...xid.ID) Option {
	return func(pm *ent.ItemMutation) {
		pm.AddLinkIDs(ids...)
	}
}

func WithDescription(v string) Option {
	return func(c *ent.ItemMutation) {
		c.SetDescription(v)
	}
}

func WithContent(v string) Option {
	return func(c *ent.ItemMutation) {
		c.SetContent(v)
	}
}

func WithVisibility(v post.Visibility) Option {
	return func(c *ent.ItemMutation) {
		c.SetVisibility(item.Visibility(v.ToEnt()))
	}
}

func WithProperties(v any) Option {
	return func(c *ent.ItemMutation) {
		c.SetProperties(v)
	}
}

func WithParentClusterAdd(id xid.ID) Option {
	return func(c *ent.ItemMutation) {
		c.AddClusterIDs(id)
	}
}

func WithParentClusterRemove(id xid.ID) Option {
	return func(c *ent.ItemMutation) {
		c.RemoveClusterIDs(id)
	}
}

func WithAssetAdd(id asset.AssetID) Option {
	return func(c *ent.ItemMutation) {
		c.AddAssetIDs(id)
	}
}

func WithAssetRemove(id asset.AssetID) Option {
	return func(c *ent.ItemMutation) {
		c.RemoveAssetIDs(id)
	}
}
