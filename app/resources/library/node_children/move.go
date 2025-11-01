package node_children

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/lexorank"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

func (d *Writer) Move(ctx context.Context, from library.QueryKey, to library.QueryKey) (*library.Node, error) {
	fromNode, err := d.nr.Get(ctx, from)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	toNode, err := d.nr.Get(ctx, to)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer tx.Rollback()

	nodes, err := tx.Node.Query().
		Select(node.FieldID).
		Where(node.ParentNodeID(xid.ID(fromNode.Mark.ID()))).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err)
	}
	childNodeIDs := dt.Map(nodes, func(c *ent.Node) xid.ID { return c.ID })

	err = tx.Node.Update().
		SetParentID(xid.ID(toNode.Mark.ID())).
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

func (d *Writer) MoveBefore(ctx context.Context, thisnode *library.Node, before library.NodeID) (*library.Node, error) {
	targetNode, siblingNode, err := d.getNodeWithSibling(ctx, before, -1)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	newsort, ok := computeNewSortKey(targetNode, siblingNode, -1)
	if !ok {
		err := d.Normalise(ctx, &targetNode.ParentNodeID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		d.logger.Warn("recalculating 'before' sort key after normalisation",
			slog.String("thisNode", thisnode.Mark.String()),
			slog.String("siblingNode", siblingNode.String()),
			slog.String("targetNode", targetNode.ID.String()),
		)

		newsort, ok = computeNewSortKey(targetNode, siblingNode, -1)
		if !ok {
			return nil, fault.New("failed to calculate new sort key after full children normalisation")
		}
	}

	err = d.db.Node.UpdateOneID(xid.ID(thisnode.Mark.ID())).SetSort(*newsort).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.nr.Get(ctx, library.NewQueryKey(thisnode.Mark))
}

func (d *Writer) MoveAfter(ctx context.Context, thisnode *library.Node, after library.NodeID) (*library.Node, error) {
	// get the sibling and the node immediately after it
	targetNode, siblingNode, err := d.getNodeWithSibling(ctx, after, 1)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	newsort, ok := computeNewSortKey(targetNode, siblingNode, 1)
	if !ok {
		err := d.Normalise(ctx, &targetNode.ParentNodeID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		d.logger.Warn("recalculating 'after' sort key after normalisation",
			slog.String("thisNode", thisnode.Mark.String()),
			slog.String("siblingNode", siblingNode.String()),
			slog.String("targetNode", targetNode.ID.String()),
		)

		newsort, ok = computeNewSortKey(targetNode, siblingNode, 1)
		if !ok {
			return nil, fault.New("failed to calculate new sort key after full children normalisation")
		}
	}

	err = d.db.Node.UpdateOneID(xid.ID(thisnode.Mark.ID())).SetSort(*newsort).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.nr.Get(ctx, library.NewQueryKey(thisnode.Mark))
}

func (d *Writer) MoveIndex(ctx context.Context, thisnode *library.Node, index int) (*library.Node, error) {
	return nil, fault.New("not implemented")
}

type reorderableNode struct {
	id   xid.ID
	sort lexorank.Key
}

func (r *reorderableNode) GetKey() lexorank.Key {
	return r.sort
}

func (r *reorderableNode) SetKey(k lexorank.Key) {
	r.sort = k
}

func (d *Writer) getNodeWithSibling(ctx context.Context, id library.NodeID, direction int) (*ent.Node, opt.Optional[*ent.Node], error) {
	targetNode, err := d.db.Node.Query().Where(node.IDEQ(xid.ID(id))).Only(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	q := d.db.Node.Query()

	if !targetNode.ParentNodeID.IsNil() {
		q.Where(node.ParentNodeID(targetNode.ParentNodeID))
	} else {
		q.Where(node.ParentNodeIDIsNil())
	}

	if direction < 0 {
		// Looking for the node before (less than current)
		q.Where(node.SortLT(targetNode.Sort)).
			Order(ent.Desc(node.FieldSort))
	} else {
		// Looking for the node after (greater than current)
		q.Where(node.SortGT(targetNode.Sort)).
			Order(ent.Asc(node.FieldSort))
	}

	siblingsResult, err := q.Limit(1).All(ctx)

	siblingNode := getSibling(siblingsResult)

	return targetNode, siblingNode, nil
}

func getSibling(siblingAfter []*ent.Node) opt.Optional[*ent.Node] {
	if len(siblingAfter) == 0 {
		return opt.NewEmpty[*ent.Node]()
	}

	return opt.New(siblingAfter[0])
}

func computeNewSortKey(targetNode *ent.Node, siblingNode opt.Optional[*ent.Node], direction int) (*lexorank.Key, bool) {
	if sibling, ok := siblingNode.Get(); ok {
		return targetNode.Sort.Between(sibling.Sort)
	} else {
		if direction < 0 {
			return targetNode.Sort.Before(1)
		} else {
			return targetNode.Sort.After(1)
		}
	}
}
