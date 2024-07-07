package collection

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type itemOptions struct {
	cid      xid.ID
	posts    []*ent.CollectionPostUpsertOne
	nodes    []*ent.CollectionNodeUpsertOne
	removals *ent.CollectionMutation
}

func WithPostAdd(id post.ID, mt MembershipType) ItemOption {
	return func(tx *ent.Tx, c *itemOptions) {
		c.posts = append(c.posts, tx.CollectionPost.Create().
			SetCollectionID(c.cid).
			SetPostID(xid.ID(id)).
			SetMembershipType(mt.String()).
			OnConflict(sql.ResolveWithNewValues()))
	}
}

func WithPostRemove(id post.ID) ItemOption {
	return func(tx *ent.Tx, c *itemOptions) {
		c.removals.RemovePostIDs(xid.ID(id))
	}
}

func WithNodeAdd(id datagraph.NodeID, mt MembershipType) ItemOption {
	return func(tx *ent.Tx, c *itemOptions) {
		c.nodes = append(c.nodes, tx.CollectionNode.Create().
			SetCollectionID(c.cid).
			SetNodeID(xid.ID(id)).
			SetMembershipType(mt.String()).
			OnConflict(sql.ResolveWithNewValues()))
	}
}

func WithNodeRemove(id datagraph.NodeID) ItemOption {
	return func(tx *ent.Tx, c *itemOptions) {
		c.removals.RemoveNodeIDs(xid.ID(id))
	}
}

func (d *database) UpdateItems(ctx context.Context, id CollectionID, opts ...ItemOption) (*Collection, error) {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer func() {
		if err != nil {
			rerr := tx.Rollback()
			if rerr != nil {
				err = fault.Wrap(
					fmt.Errorf("rollback failed: %w: original error: %w", rerr, err),
					fctx.With(ctx),
				)
			}
		}
	}()

	removals := tx.Collection.UpdateOneID(xid.ID(id))

	options := itemOptions{
		cid:      xid.ID(id),
		posts:    []*ent.CollectionPostUpsertOne{},
		nodes:    []*ent.CollectionNodeUpsertOne{},
		removals: removals.Mutation(),
	}

	for _, fn := range opts {
		fn(tx, &options)
	}

	if len(removals.Mutation().RemovedEdges()) > 0 {
		err := removals.Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	for _, p := range options.posts {
		err = p.Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	for _, n := range options.nodes {
		err = n.Exec(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, id)
}
