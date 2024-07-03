package collection

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type (
	Option func(*ent.CollectionMutation)
	Filter func(*ent.CollectionQuery)
)

type Repository interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		desc string,
		opts ...Option) (*Collection, error)

	List(ctx context.Context, filters ...Filter) ([]*Collection, error)
	Get(ctx context.Context, id CollectionID) (*Collection, error)

	Update(ctx context.Context, id CollectionID, opts ...Option) (*Collection, error)

	Delete(ctx context.Context, id CollectionID) error
}

func WithID(id CollectionID) Option {
	return func(c *ent.CollectionMutation) {
		c.SetID(xid.ID(id))
	}
}

func WithName(v string) Option {
	return func(c *ent.CollectionMutation) {
		c.SetName(v)
	}
}

func WithDescription(v string) Option {
	return func(c *ent.CollectionMutation) {
		c.SetDescription(v)
	}
}

func WithPostAdd(id post.ID) Option {
	return func(c *ent.CollectionMutation) {
		c.AddPostIDs(xid.ID(id))
	}
}

func WithPostRemove(id post.ID) Option {
	return func(c *ent.CollectionMutation) {
		c.RemovePostIDs(xid.ID(id))
	}
}

func WithNodeAdd(id datagraph.NodeID) Option {
	return func(c *ent.CollectionMutation) {
		c.AddNodeIDs(xid.ID(id))
	}
}

func WithNodeRemove(id datagraph.NodeID) Option {
	return func(c *ent.CollectionMutation) {
		c.RemoveNodeIDs(xid.ID(id))
	}
}
