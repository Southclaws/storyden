package collection_item

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/collection/collection_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/collectionnode"
	"github.com/Southclaws/storyden/internal/ent/collectionpost"
)

type Repository struct {
	db      *ent.Client
	querier *collection_querier.Querier
	roles   *role_repo.Repository
}

func New(db *ent.Client, querier *collection_querier.Querier, roles *role_repo.Repository) *Repository {
	return &Repository{db: db, querier: querier, roles: roles}
}

type itemChange struct {
	id     xid.ID
	mt     collection.MembershipType
	t      datagraph.Kind
	remove bool
}

type itemChanges []itemChange

type ItemOption func(*itemChanges)

func WithPost(id post.ID, mt collection.MembershipType) ItemOption {
	return func(c *itemChanges) {
		*c = append(*c, itemChange{
			id: xid.ID(id),
			mt: mt,
			t:  datagraph.KindPost,
		})
	}
}

func WithPostRemove(id post.ID) ItemOption {
	return func(c *itemChanges) {
		*c = append(*c, itemChange{
			id:     xid.ID(id),
			t:      datagraph.KindPost,
			remove: true,
		})
	}
}

func WithNode(id library.NodeID, mt collection.MembershipType) ItemOption {
	return func(c *itemChanges) {
		*c = append(*c, itemChange{
			id: xid.ID(id),
			mt: mt,
			t:  datagraph.KindNode,
		})
	}
}

func WithNodeRemove(id library.NodeID) ItemOption {
	return func(c *itemChanges) {
		*c = append(*c, itemChange{
			id:     xid.ID(id),
			t:      datagraph.KindNode,
			remove: true,
		})
	}
}

func (d *Repository) UpdateItems(ctx context.Context, qk collection.QueryKey, opts ...ItemOption) (*collection.CollectionWithItems, error) {
	cid, err := d.db.Collection.Query().Where(qk.Predicate()).OnlyID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	options := itemChanges{}

	for _, fn := range opts {
		fn(&options)
	}

	for _, op := range options {
		var err error

		switch op.t {
		case datagraph.KindPost:
			predicate := collectionpost.And(
				collectionpost.CollectionID(cid),
				collectionpost.PostID(op.id),
			)

			if op.remove {
				_, err = d.db.CollectionPost.Delete().Where(predicate).Exec(ctx)
			} else {
				exists, exerr := d.db.CollectionPost.Query().Where(predicate).Exist(ctx)
				if exerr != nil {
					return nil, exerr
				}

				if exists {
					err = d.db.CollectionPost.Update().
						Where(predicate).
						SetMembershipType(op.mt.String()).
						Exec(ctx)
				} else {
					err = d.db.CollectionPost.Create().
						SetCollectionID(cid).
						SetPostID(op.id).
						SetMembershipType(op.mt.String()).
						Exec(ctx)
				}
			}

		case datagraph.KindNode:
			predicate := collectionnode.And(
				collectionnode.CollectionID(cid),
				collectionnode.NodeID(op.id),
			)

			if op.remove {
				_, err = d.db.CollectionNode.Delete().Where(
					collectionnode.CollectionID(cid),
					collectionnode.NodeID(op.id),
				).Exec(ctx)
			} else {
				exists, exerr := d.db.CollectionNode.Query().Where(predicate).Exist(ctx)
				if exerr != nil {
					return nil, exerr
				}

				if exists {
					err = d.db.CollectionNode.Update().
						Where(predicate).
						SetMembershipType(op.mt.String()).
						Exec(ctx)
				} else {
					err = d.db.CollectionNode.Create().
						SetCollectionID(cid).
						SetNodeID(op.id).
						SetMembershipType(op.mt.String()).
						Exec(ctx)
				}

			}
		}
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return d.querier.Get(ctx, qk)
}

func (d *Repository) ProbeItem(ctx context.Context, qk collection.QueryKey, itemID xid.ID) (*collection.CollectionItemStatus, error) {
	r, err := d.db.Collection.Query().
		Where(qk.Predicate()).
		WithOwner().
		WithCollectionPosts(func(cnq *ent.CollectionPostQuery) {
			cnq.Where(collectionpost.PostID(itemID)).
				WithPost(func(pq *ent.PostQuery) {
					pq.WithAuthor()
				})
		}).
		WithCollectionNodes(func(cnq *ent.CollectionNodeQuery) {
			cnq.Where(collectionnode.NodeID(itemID)).
				WithNode(func(nq *ent.NodeQuery) {
					nq.WithOwner()
				})
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accs := []*ent.Account{r.Edges.Owner}
	if len(r.Edges.CollectionNodes) > 0 {
		if cn := r.Edges.CollectionNodes[0]; cn != nil && cn.Edges.Node != nil && cn.Edges.Node.Edges.Owner != nil {
			accs = append(accs, cn.Edges.Node.Edges.Owner)
		}
	}
	if len(r.Edges.CollectionPosts) > 0 {
		if cp := r.Edges.CollectionPosts[0]; cp != nil && cp.Edges.Post != nil && cp.Edges.Post.Edges.Author != nil {
			accs = append(accs, cp.Edges.Post.Edges.Author)
		}
	}

	filteredAccs := make([]*ent.Account, 0, len(accs))
	for _, acc := range accs {
		if acc == nil {
			continue
		}

		filteredAccs = append(filteredAccs, acc)
	}

	roleHydrator, err := d.roles.BuildMultiHydrator(ctx, filteredAccs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	item, err := func() (_ opt.Optional[collection.CollectionItem], err error) {
		var i *collection.CollectionItem
		if len(r.Edges.CollectionNodes) > 0 {
			i, err = collection.MapNode(r.Edges.CollectionNodes[0], roleHydrator.Hydrate)
		}
		if len(r.Edges.CollectionPosts) > 0 {
			i, err = collection.MapPost(r.Edges.CollectionPosts[0], roleHydrator.Hydrate)
		}
		return opt.NewPtr(i), err
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	col, err := collection.Map(nil, roleHydrator.Hydrate)(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &collection.CollectionItemStatus{
		Collection: *col,
		Item:       item,
	}, nil
}
