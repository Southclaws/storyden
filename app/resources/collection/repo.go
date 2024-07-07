package collection

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/collectionnode"
	"github.com/Southclaws/storyden/internal/ent/collectionpost"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type (
	Option     func(*ent.CollectionMutation)
	ItemOption func(*itemChanges)
	Filter     func(*ent.CollectionQuery)
	ItemFilter func(*ent.CollectionPostQuery, *ent.CollectionNodeQuery)
)

type Repository interface {
	Create(ctx context.Context,
		owner account.AccountID,
		name string,
		desc string,
		opts ...Option) (*CollectionWithItems, error)

	List(ctx context.Context, filters ...Filter) ([]*Collection, error)
	Get(ctx context.Context, id CollectionID, filters ...ItemFilter) (*CollectionWithItems, error)

	ProbeItem(ctx context.Context, id CollectionID, item xid.ID) (*CollectionItemStatus, error)

	Update(ctx context.Context, id CollectionID, opts ...Option) (*CollectionWithItems, error)
	UpdateItems(ctx context.Context, id CollectionID, opts ...ItemOption) (*CollectionWithItems, error)

	Delete(ctx context.Context, id CollectionID) error
}

func WithVisibility(v ...visibility.Visibility) ItemFilter {
	return func(pq *ent.CollectionPostQuery, nq *ent.CollectionNodeQuery) {
		pv := dt.Map(v, func(v visibility.Visibility) ent_post.Visibility { return ent_post.Visibility(v.String()) })
		pq.Where(
			collectionpost.HasPostWith(
				ent_post.VisibilityIn(pv...),
			),
		)

		nv := dt.Map(v, func(v visibility.Visibility) ent_node.Visibility { return ent_node.Visibility(v.String()) })
		nq.Where(
			collectionnode.HasNodeWith(
				ent_node.VisibilityIn(nv...),
			),
		)
	}
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
