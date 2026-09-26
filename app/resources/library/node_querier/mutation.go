package node_querier

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/nodeversion"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type MutationState struct {
	Mark       library.Mark
	Name       string
	Content    opt.Optional[datagraph.Content]
	Visibility visibility.Visibility
	UpdatedAt  time.Time
	HasDraft   bool
}

// MutationStates reads the state needed to prepare bulk page mutations without
// hydrating complete pages and their unrelated edges. Content supports stable
// block IDs; visibility, timestamps, and draft records support mutation guards.
// Results are keyed by node ID; missing pages are omitted. Explicit slugs are
// matched as slugs even when they resemble IDs. This method performs no writes.
func (q *Querier) MutationStates(ctx context.Context, keys []library.QueryKey, slugs ...mark.Slug) (map[library.NodeID]MutationState, error) {
	states := make(map[library.NodeID]MutationState, len(keys))
	if len(keys) == 0 && len(slugs) == 0 {
		return states, nil
	}

	predicates := make([]predicate.Node, len(keys))
	for i, key := range keys {
		predicates[i] = key.Predicate()
	}

	for _, slug := range slugs {
		predicates = append(predicates, node.Slug(slug.String()))
	}

	rows, err := q.db.Node.Query().Where(node.Or(predicates...)).
		Select(node.FieldID, node.FieldSlug, node.FieldName, node.FieldContent, node.FieldVisibility, node.FieldUpdatedAt).
		WithVersions(func(v *ent.NodeVersionQuery) {
			v.Where(nodeversion.StatusEQ(nodeversion.StatusDraft)).Select(nodeversion.FieldID, nodeversion.FieldNodeID)
		}).All(ctx)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		content := opt.NewEmpty[datagraph.Content]()
		if row.Content != nil {
			parsed, err := datagraph.NewRichTextWithBlocks(*row.Content)
			if err != nil {
				return nil, err
			}

			content = opt.New(parsed.Content)
		}

		v, err := visibility.NewVisibility(string(row.Visibility))
		if err != nil {
			return nil, err
		}

		states[library.NodeID(row.ID)] = MutationState{
			Mark:       library.NewMark(row.ID, row.Slug),
			Name:       row.Name,
			Content:    content,
			Visibility: v,
			UpdatedAt:  row.UpdatedAt,
			HasDraft:   len(row.Edges.Versions) != 0,
		}
	}

	return states, nil
}

// Ancestors returns the parent relation for the requested nodes and every
// ancestor reachable from them.
func (q *Querier) Ancestors(ctx context.Context, ids []library.NodeID) (map[library.NodeID]library.NodeID, error) {
	parents := make(map[library.NodeID]library.NodeID)
	rows, err := q.queryAncestry(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row.ParentID != nil {
			parents[library.NodeID(row.ID)] = library.NodeID(*row.ParentID)
		}
	}

	return parents, nil
}

// Ancestry returns the ancestors of a node ordered from the root to its parent.
// The node itself is not included. Visibility rules omit the first inaccessible
// ancestor and the remainder of the path above it.
func (q *Querier) Ancestry(ctx context.Context, id library.NodeID, opts ...Option) ([]library.NodeReference, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	rows, err := q.queryAncestry(ctx, []library.NodeID{id})
	if err != nil {
		return nil, err
	}

	parents := make(map[library.NodeID]library.NodeID, len(rows))
	byID := make(map[library.NodeID]ancestryRow, len(rows))
	for _, row := range rows {
		nodeID := library.NodeID(row.ID)
		byID[nodeID] = row

		if row.ParentID != nil {
			parents[nodeID] = library.NodeID(*row.ParentID)
		}
	}

	ancestry := make([]library.NodeReference, 0, len(rows))
	seen := map[library.NodeID]bool{id: true}
	for ancestorID, ok := parents[id]; ok; ancestorID, ok = parents[ancestorID] {
		if seen[ancestorID] {
			return nil, fmt.Errorf("cycle in node ancestry at %s", ancestorID)
		}
		seen[ancestorID] = true

		row, ok := byID[ancestorID]
		if !ok {
			return nil, fmt.Errorf("missing node ancestry record for %s", ancestorID)
		}
		if !canViewAncestryRow(o, row) {
			break
		}

		ancestry = append(ancestry, library.NodeReference{
			Mark: library.NewMark(row.ID, row.Slug),
			Name: row.Name,
		})
	}

	slices.Reverse(ancestry)

	return ancestry, nil
}

type ancestryRow struct {
	ID         xid.ID          `db:"id"`
	ParentID   *xid.ID         `db:"parent_node_id"`
	AccountID  xid.ID          `db:"account_id"`
	Name       string          `db:"name"`
	Slug       string          `db:"slug"`
	Visibility node.Visibility `db:"visibility"`
}

func canViewAncestryRow(o *options, row ancestryRow) bool {
	if !o.visibilityRules {
		return true
	}

	acc, ok := o.requestingAccount.Get()
	if !ok {
		return row.Visibility == node.VisibilityPublished
	}
	if row.AccountID == xid.ID(acc.ID) || row.Visibility == node.VisibilityPublished {
		return true
	}

	canViewInReview := acc.Roles.Permissions().HasAny(rbac.PermissionAdministrator, rbac.PermissionManageLibrary)

	return canViewInReview && row.Visibility == node.VisibilityReview
}

func (q *Querier) queryAncestry(ctx context.Context, ids []library.NodeID) ([]ancestryRow, error) {
	if len(ids) == 0 {
		return []ancestryRow{}, nil
	}

	values := make([]string, len(ids))
	for i, id := range ids {
		values[i] = id.String()
	}

	query, args, err := sqlx.In(`WITH RECURSIVE ancestry AS (
  SELECT id, parent_node_id, account_id, name, slug, visibility FROM nodes WHERE id IN (?)
  UNION
  SELECT n.id, n.parent_node_id, n.account_id, n.name, n.slug, n.visibility FROM nodes n JOIN ancestry a ON n.id = a.parent_node_id
) SELECT id, parent_node_id, account_id, name, slug, visibility FROM ancestry`, values)
	if err != nil {
		return nil, err
	}

	rows := []ancestryRow{}
	if err := q.raw.SelectContext(ctx, &rows, q.raw.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("read node ancestry: %w", err)
	}

	return rows, nil
}
