package node_writer

import (
	"context"

	"github.com/Southclaws/lexorank"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

func setBatchSortKeys(ctx context.Context, db *ent.Client, builders []*ent.NodeCreate) error {
	groups := make(map[xid.ID][]*ent.NodeCreate)
	for _, builder := range builders {
		parent, _ := builder.Mutation().ParentID()
		groups[parent] = append(groups[parent], builder)
	}

	parents := make([]predicate.Node, 0, len(groups))
	for parent := range groups {
		parents = append(parents, siblingPredicate(parent))
	}

	var maxima []struct {
		ParentNodeID *xid.ID `sql:"parent_node_id"`
		Max          lexorank.Key
	}
	err := db.Node.Query().
		Where(node.Or(parents...)).
		GroupBy(node.FieldParentNodeID).
		Aggregate(ent.Max(node.FieldSort)).
		Scan(ctx, &maxima)
	if err != nil {
		return err
	}

	last := make(map[xid.ID]lexorank.Key, len(maxima))
	for _, row := range maxima {
		parent := xid.NilID()
		if row.ParentNodeID != nil {
			parent = *row.ParentNodeID
		}

		last[parent] = row.Max
	}

	exhausted := make(map[xid.ID][]*ent.NodeCreate)
	for parent, children := range groups {
		key, exists := last[parent]
		for _, child := range children {
			next := &lexorank.Middle
			if exists {
				var ok bool
				next, ok = key.After(100)
				if !ok {
					exhausted[parent] = children
					break
				}
			}

			child.SetSort(*next)
			key, exists = *next, true
		}
	}

	if len(exhausted) == 0 {
		return nil
	}

	return rebalanceBatchSiblings(ctx, db, exhausted)
}

func siblingPredicate(parent xid.ID) predicate.Node {
	if parent == xid.NilID() {
		return node.ParentNodeIDIsNil()
	}

	return node.ParentNodeID(parent)
}

func rebalanceBatchSiblings(ctx context.Context, db *ent.Client, groups map[xid.ID][]*ent.NodeCreate) error {
	parents := make([]predicate.Node, 0, len(groups))
	for parent := range groups {
		parents = append(parents, siblingPredicate(parent))
	}

	rows, err := db.Node.Query().
		Where(node.Or(parents...)).
		Select(node.FieldID, node.FieldParentNodeID, node.FieldSort).
		Order(ent.Asc(node.FieldSort), ent.Asc(node.FieldID)).
		All(ctx)
	if err != nil {
		return err
	}

	siblings := make(map[xid.ID][]*ent.Node)
	for _, row := range rows {
		siblings[row.ParentNodeID] = append(siblings[row.ParentNodeID], row)
	}

	mutations := make(map[xid.ID]*ent.NodeMutation, len(rows))
	ids := make([]xid.ID, 0, len(rows))
	for parent, children := range groups {
		existing := siblings[parent]
		total := len(existing) + len(children)
		for i, row := range existing {
			mutation := db.Node.Update().Mutation()
			mutation.SetSort(lexorank.KeyAt(0, float64(i+2)/float64(total+3)))
			mutations[row.ID] = mutation
			ids = append(ids, row.ID)
		}

		for i, child := range children {
			child.SetSort(lexorank.KeyAt(0, float64(len(existing)+i+2)/float64(total+3)))
		}
	}

	return db.Node.Update().Where(node.IDIn(ids...)).Modify(batchUpdate(mutations)).Exec(ctx)
}
