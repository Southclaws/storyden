package collection

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/collectionnode"
	"github.com/Southclaws/storyden/internal/ent/collectionpost"
)

type itemChange struct {
	id     xid.ID
	mt     MembershipType
	t      datagraph.Kind
	remove bool
}

type itemChanges []itemChange

func WithPost(id post.ID, mt MembershipType) ItemOption {
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

func WithNode(id datagraph.NodeID, mt MembershipType) ItemOption {
	return func(c *itemChanges) {
		*c = append(*c, itemChange{
			id: xid.ID(id),
			mt: mt,
			t:  datagraph.KindNode,
		})
	}
}

func WithNodeRemove(id datagraph.NodeID) ItemOption {
	return func(c *itemChanges) {
		*c = append(*c, itemChange{
			id:     xid.ID(id),
			t:      datagraph.KindNode,
			remove: true,
		})
	}
}

func (d *database) UpdateItems(ctx context.Context, id CollectionID, opts ...ItemOption) (*CollectionWithItems, error) {
	options := itemChanges{}

	for _, fn := range opts {
		fn(&options)
	}

	cid := xid.ID(id)

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
					collectionnode.CollectionID(xid.ID(id)),
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

	return d.Get(ctx, id)
}

func (d *database) ProbeItem(ctx context.Context, id CollectionID, itemID xid.ID) (*CollectionItemStatus, error) {
	r, err := d.db.Collection.Query().
		Where(collection.ID(xid.ID(id))).
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

	item, err := func() (_ opt.Optional[CollectionItem], err error) {
		var i *CollectionItem
		if len(r.Edges.CollectionNodes) > 0 {
			i, err = MapCollectionNode(r.Edges.CollectionNodes[0])
		}
		if len(r.Edges.CollectionPosts) > 0 {
			i, err = MapCollectionPost(r.Edges.CollectionPosts[0])
		}
		return opt.NewPtr(i), err
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	collection, err := MapCollection(r)
	if err != nil {
		return nil, err
	}

	return &CollectionItemStatus{
		Collection: *collection,
		Item:       item,
	}, nil
}
