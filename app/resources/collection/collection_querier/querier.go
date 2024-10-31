package collection_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_collection "github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/collectionnode"
	"github.com/Southclaws/storyden/internal/ent/collectionpost"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db}
}

type listOption struct {
	ownerID               opt.Optional[string]
	queryForIncludedItems opt.Optional[[]xid.ID]
}

type (
	ItemFilter func(*ent.CollectionPostQuery, *ent.CollectionNodeQuery)
	Option     func(*listOption)
)

func WithOwnerHandle(v string) Option {
	return func(c *listOption) {
		c.ownerID = opt.New(v)
	}
}

func WithItemPresenceQuery(id xid.ID) Option {
	return func(lo *listOption) {
		lo.queryForIncludedItems = opt.New([]xid.ID{id})
	}
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

func (d *Querier) List(ctx context.Context, filters ...Option) ([]*collection.Collection, error) {
	var opts listOption
	for _, fn := range filters {
		fn(&opts)
	}

	q := d.db.Collection.
		Query().
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithCollectionPosts().
		WithCollectionNodes()

	opts.ownerID.Call(func(v string) {
		q.Where(ent_collection.HasOwnerWith(ent_account.Handle(v)))
	})

	cols, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all, err := collection.MapList(opts.queryForIncludedItems.OrZero(), cols)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return all, nil
}

func (d *Querier) Get(ctx context.Context, id collection.CollectionID, filters ...ItemFilter) (*collection.CollectionWithItems, error) {
	filters = append(filters, func(pcq *ent.CollectionPostQuery, ncq *ent.CollectionNodeQuery) {
		if pcq != nil {
			pcq.WithPost(func(pq *ent.PostQuery) {
				pq.WithAuthor(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				})
				pq.WithCategory()
				pq.WithTags()
				pq.WithRoot()
			})
		}

		if ncq != nil {
			ncq.WithNode(func(nq *ent.NodeQuery) {
				nq.WithOwner(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				})
				nq.WithAssets()
				nq.WithTags()
			})
		}
	})

	col, err := d.db.Collection.
		Query().
		Where(ent_collection.ID(xid.ID(id))).
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithCollectionPosts(func(pq *ent.CollectionPostQuery) {
			for _, fn := range filters {
				fn(pq, nil)
			}
		}).
		WithCollectionNodes(func(nq *ent.CollectionNodeQuery) {
			for _, fn := range filters {
				fn(nil, nq)
			}
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return collection.MapWithItems(col)
}

func (d *Querier) Probe(ctx context.Context, id collection.CollectionID) (*collection.Collection, error) {
	col, err := d.db.Collection.
		Query().
		Where(ent_collection.ID(xid.ID(id))).
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return collection.Map(nil)(col)
}
