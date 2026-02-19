package collection_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
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
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{
		db:          db,
		roleQuerier: roleQuerier,
	}
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
		WithOwner().
		WithCollectionPosts().
		WithCollectionNodes()

	opts.ownerID.Call(func(v string) {
		q.Where(ent_collection.HasOwnerWith(ent_account.Handle(v)))
	})

	cols, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accs := dt.Map(cols, func(col *ent.Collection) *ent.Account {
		return col.Edges.Owner
	})

	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, accs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	all, err := collection.MapList(opts.queryForIncludedItems.OrZero(), cols, roleHydrator.Hydrate)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return all, nil
}

func (d *Querier) Get(ctx context.Context, qk collection.QueryKey, filters ...ItemFilter) (*collection.CollectionWithItems, error) {
	filters = append(filters, func(pcq *ent.CollectionPostQuery, ncq *ent.CollectionNodeQuery) {
		if pcq != nil {
			pcq.WithPost(func(pq *ent.PostQuery) {
				pq.WithAuthor()
				pq.WithCategory()
				pq.WithTags()
				pq.WithRoot()
			})
		}

		if ncq != nil {
			ncq.WithNode(func(nq *ent.NodeQuery) {
				nq.WithOwner()
				nq.WithAssets()
				nq.WithTags()
			})
		}
	})

	col, err := d.db.Collection.
		Query().
		Where(qk.Predicate()).
		WithOwner().
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

	accs := roleHydrationTargetsFromCollection(col)
	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, accs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return collection.MapWithItems(col, roleHydrator.Hydrate)
}

func (d *Querier) Probe(ctx context.Context, qk collection.QueryKey) (*collection.Collection, error) {
	col, err := d.db.Collection.
		Query().
		Where(qk.Predicate()).
		WithOwner().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, []*ent.Account{col.Edges.Owner})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return collection.Map(nil, roleHydrator.Hydrate)(col)
}

func roleHydrationTargetsFromCollection(col *ent.Collection) []*ent.Account {
	targets := map[xid.ID]*ent.Account{}

	if owner := col.Edges.Owner; owner != nil {
		targets[owner.ID] = owner
	}

	for _, cp := range col.Edges.CollectionPosts {
		if cp == nil || cp.Edges.Post == nil || cp.Edges.Post.Edges.Author == nil {
			continue
		}

		author := cp.Edges.Post.Edges.Author
		targets[author.ID] = author
	}

	for _, cn := range col.Edges.CollectionNodes {
		if cn == nil || cn.Edges.Node == nil || cn.Edges.Node.Edges.Owner == nil {
			continue
		}

		owner := cn.Edges.Node.Edges.Owner
		targets[owner.ID] = owner
	}

	out := make([]*ent.Account, 0, len(targets))
	for _, acc := range targets {
		out = append(out, acc)
	}

	return out
}
