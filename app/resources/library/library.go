package library

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type (
	Option func(*ent.NodeMutation)
	Filter func(*ent.NodeQuery)
)

type Repository interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		opts ...Option,
	) (*Node, error)

	Get(ctx context.Context, slug NodeSlug) (*Node, error)
	GetByID(ctx context.Context, id NodeID) (*Node, error)

	// Update a node by ID.
	// NOTE: slug based update is not supported at the repo level because you'll
	// probably always have a node ID in context anyway and it makes changing
	// the actual slug a bit more complex due to the naïve implementation.
	Update(ctx context.Context, id NodeID, opts ...Option) (*Node, error)

	// Delete removes a node permanently, it does not manage children.
	Delete(ctx context.Context, slug NodeSlug) error
}

func WithID(id NodeID) Option {
	return func(c *ent.NodeMutation) {
		c.SetID(xid.ID(id))
	}
}

func WithName(v string) Option {
	return func(c *ent.NodeMutation) {
		c.SetName(v)
	}
}

func WithSlug(v string) Option {
	return func(c *ent.NodeMutation) {
		c.SetSlug(v)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.NodeMutation) {
		m.AddAssetIDs(a...)
	}
}

func WithAssetsRemoved(a []asset.AssetID) Option {
	return func(m *ent.NodeMutation) {
		m.RemoveAssetIDs(a...)
	}
}

func WithLink(id xid.ID) Option {
	return func(pm *ent.NodeMutation) {
		pm.SetLinkID(id)
	}
}

func WithContentLinks(ids ...xid.ID) Option {
	return func(pm *ent.NodeMutation) {
		pm.AddContentLinkIDs(ids...)
	}
}

func WithContent(v content.Rich) Option {
	return func(c *ent.NodeMutation) {
		c.SetContent(v.HTML())
		c.SetDescription(v.Short())
	}
}

func WithDescription(v string) Option {
	return func(c *ent.NodeMutation) {
		c.SetDescription(v)
	}
}

func WithParent(v NodeID) Option {
	return func(c *ent.NodeMutation) {
		c.SetParentID(xid.ID(v))
	}
}

func WithVisibility(v visibility.Visibility) Option {
	return func(c *ent.NodeMutation) {
		c.SetVisibility(node.Visibility(v.String()))
	}
}

func WithMetadata(v map[string]any) Option {
	return func(c *ent.NodeMutation) {
		c.SetMetadata(v)
	}
}

func WithChildNodeAdd(id xid.ID) Option {
	return func(c *ent.NodeMutation) {
		c.AddNodeIDs(id)
	}
}

func WithChildNodeRemove(id xid.ID) Option {
	return func(c *ent.NodeMutation) {
		c.RemoveNodeIDs(id)
	}
}
