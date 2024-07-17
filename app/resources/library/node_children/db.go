package node_children

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type database struct {
	db *ent.Client
	nr library.Repository
}

func New(db *ent.Client, nr library.Repository) Repository {
	return &database{db, nr}
}

func (d *database) Move(ctx context.Context, fromSlug library.NodeSlug, toSlug library.NodeSlug) (*library.Node, error) {
	fromNode, err := d.nr.Get(ctx, fromSlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	toNode, err := d.nr.Get(ctx, toSlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := d.db.Node.Query().
		Select(node.FieldID).
		Where(node.ParentNodeID(xid.ID(fromNode.ID))).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err)
	}
	childNodeIDs := dt.Map(nodes, func(c *ent.Node) xid.ID { return c.ID })

	err = d.db.Node.Update().
		SetParentID(xid.ID(toNode.ID)).
		Where(node.IDIn(childNodeIDs...)).
		Exec(ctx)
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			panic(fmt.Errorf("while handling error: %w, rollback error: %s", err, terr))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return toNode, nil
}
