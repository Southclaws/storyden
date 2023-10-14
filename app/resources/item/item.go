package item

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type (
	Option func(*ent.ItemMutation)
	Filter func(*ent.ItemQuery)
)

type ItemRepository interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		desc string,
		opts ...Option,
	) (*datagraph.Item, error)

	List(ctx context.Context, filters ...Filter) ([]*datagraph.Item, error)
	Get(ctx context.Context, slug datagraph.ItemSlug) (*datagraph.Item, error)

	Update(ctx context.Context, slug datagraph.ItemSlug, opts ...Option) (*datagraph.Item, error)

	Delete(ctx context.Context, slug datagraph.ItemSlug) error
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

func WithImageURL(v string) Option {
	return func(c *ent.ItemMutation) {
		c.SetImageURL(v)
	}
}

func WithDescription(v string) Option {
	return func(c *ent.ItemMutation) {
		c.SetDescription(v)
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

func WithAssetAdd(id string) Option {
	return func(c *ent.ItemMutation) {
		c.AddAssetIDs(id)
	}
}

func WithAssetRemove(id string) Option {
	return func(c *ent.ItemMutation) {
		c.RemoveAssetIDs(id)
	}
}
