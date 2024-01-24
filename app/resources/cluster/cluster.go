package cluster

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type (
	Option func(*ent.ClusterMutation)
	Filter func(*ent.ClusterQuery)
)

type Repository interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		slug string,
		desc string,
		opts ...Option,
	) (*datagraph.Cluster, error)

	Get(ctx context.Context, slug datagraph.ClusterSlug) (*datagraph.Cluster, error)

	// Update a cluster by ID.
	// NOTE: slug based update is not supported at the repo level because you'll
	// probably always have a cluster ID in context anyway and it makes changing
	// the actual slug a bit more complex due to the naïve implementation.
	Update(ctx context.Context, id datagraph.ClusterID, opts ...Option) (*datagraph.Cluster, error)

	// Delete removes a cluster permanently, it does not manage children.
	Delete(ctx context.Context, slug datagraph.ClusterSlug) error
}

func WithID(id datagraph.ClusterID) Option {
	return func(c *ent.ClusterMutation) {
		c.SetID(xid.ID(id))
	}
}

func WithName(v string) Option {
	return func(c *ent.ClusterMutation) {
		c.SetName(v)
	}
}

func WithSlug(v string) Option {
	return func(c *ent.ClusterMutation) {
		c.SetSlug(v)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.ClusterMutation) {
		m.AddAssetIDs(dt.Map(a, func(id asset.AssetID) string { return string(id) })...)
	}
}

func WithAssetsRemoved(a []asset.AssetID) Option {
	return func(m *ent.ClusterMutation) {
		m.RemoveAssetIDs(dt.Map(a, func(id asset.AssetID) string { return string(id) })...)
	}
}

func WithLinks(ids ...xid.ID) Option {
	return func(pm *ent.ClusterMutation) {
		pm.AddLinkIDs(ids...)
	}
}

func WithDescription(v string) Option {
	return func(c *ent.ClusterMutation) {
		c.SetDescription(v)
	}
}

func WithContent(v string) Option {
	return func(c *ent.ClusterMutation) {
		c.SetContent(v)
	}
}

func WithProperties(v any) Option {
	return func(c *ent.ClusterMutation) {
		c.SetProperties(v)
	}
}

func WithChildClusterAdd(id xid.ID) Option {
	return func(c *ent.ClusterMutation) {
		c.AddClusterIDs(id)
	}
}

func WithChildClusterRemove(id xid.ID) Option {
	return func(c *ent.ClusterMutation) {
		c.RemoveClusterIDs(id)
	}
}

func WithItemAdd(id xid.ID) Option {
	return func(c *ent.ClusterMutation) {
		c.AddItemIDs(id)
	}
}

func WithItemRemove(id xid.ID) Option {
	return func(c *ent.ClusterMutation) {
		c.RemoveItemIDs(id)
	}
}
