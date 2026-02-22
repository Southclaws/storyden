package collection_item

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
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
	db          *ent.Client
	querier     *collection_querier.Querier
	roleQuerier *role_querier.Querier
}

func New(db *ent.Client, querier *collection_querier.Querier, roleQuerier *role_querier.Querier) *Repository {
	return &Repository{db: db, querier: querier, roleQuerier: roleQuerier}
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

	roleTargets := make([]*ent.Account, 0, 3)
	if owner := r.Edges.Owner; owner != nil {
		roleTargets = append(roleTargets, owner)
	}
	if len(r.Edges.CollectionNodes) > 0 && r.Edges.CollectionNodes[0] != nil {
		if nodeEdge := r.Edges.CollectionNodes[0].Edges.Node; nodeEdge != nil && nodeEdge.Edges.Owner != nil {
			roleTargets = append(roleTargets, nodeEdge.Edges.Owner)
		}
	}
	if len(r.Edges.CollectionPosts) > 0 && r.Edges.CollectionPosts[0] != nil {
		if postEdge := r.Edges.CollectionPosts[0].Edges.Post; postEdge != nil && postEdge.Edges.Author != nil {
			roleTargets = append(roleTargets, postEdge.Edges.Author)
		}
	}
	if err := d.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	item, err := func() (_ opt.Optional[collection.CollectionItem], err error) {
		var i *collection.CollectionItem
		if len(r.Edges.CollectionNodes) > 0 {
			i, err = collection.MapNode(r.Edges.CollectionNodes[0])
		}
		if len(r.Edges.CollectionPosts) > 0 {
			i, err = collection.MapPost(r.Edges.CollectionPosts[0])
		}
		return opt.NewPtr(i), err
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	col, err := collection.Map(nil)(r)
	if err != nil {
		return nil, err
	}

	return &collection.CollectionItemStatus{
		Collection: *col,
		Item:       item,
	}, nil
}
