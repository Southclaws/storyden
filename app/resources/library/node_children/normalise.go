package node_children

import (
	"context"
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/lexorank"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

func (d *Writer) Normalise(ctx context.Context, parent *xid.ID) error {
	d.logger.Warn("normalising node sort order", slog.String("parent", opt.NewPtr(parent).String()))

	q := d.db.Node.Query().Order(ent.Asc(node.FieldSort))

	if parent != nil {
		q.Where(node.ParentNodeID(*parent))
	} else {
		q.Where(node.ParentNodeIDIsNil())
	}

	nodes, err := q.All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if len(nodes) == 0 {
		return nil
	}

	rol := lexorank.ReorderableList(dt.Map(nodes, func(n *ent.Node) lexorank.Reorderable {
		return &reorderableNode{id: n.ID, sort: n.Sort}
	}))

	rol.Normalise()

	// NOTE: This transaction deadlocked a bit in the tests so I think on larger
	// nodes with many children during high traffic this could be a problem.
	// TODO: Explore potential background job alternatives for full normalise
	// and do a smaller rebalance-until-safe for insertion time rebalance.
	// tx, err := d.db.Tx(ctx)
	// if err != nil {
	// 	return fault.Wrap(err, fctx.With(ctx))
	// }
	// defer tx.Rollback()

	for _, l := range rol {
		n := l.(*reorderableNode)

		err := d.db.Node.UpdateOneID(n.id).
			SetSort(n.sort).
			Exec(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	// if err := tx.Commit(); err != nil {
	// 	return fault.Wrap(err, fctx.With(ctx))
	// }

	return nil
}
