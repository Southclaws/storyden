package node_querier

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mark"
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

func (q *Querier) Ancestors(ctx context.Context, ids []library.NodeID) (map[library.NodeID]library.NodeID, error) {
	parents := make(map[library.NodeID]library.NodeID)
	if len(ids) == 0 {
		return parents, nil
	}

	values := make([]string, len(ids))
	for i, id := range ids {
		values[i] = id.String()
	}

	query, args, err := sqlx.In(`WITH RECURSIVE ancestry AS (
  SELECT id, parent_node_id FROM nodes WHERE id IN (?)
  UNION
  SELECT n.id, n.parent_node_id FROM nodes n JOIN ancestry a ON n.id = a.parent_node_id
	) SELECT id, parent_node_id FROM ancestry`, values)
	if err != nil {
		return nil, err
	}

	rows, err := q.db.QueryContext(ctx, q.raw.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id xid.ID
		var parent *xid.ID
		if err := rows.Scan(&id, &parent); err != nil {
			return nil, fmt.Errorf("read node ancestry: %w", err)
		}

		if parent != nil {
			parents[library.NodeID(id)] = library.NodeID(*parent)
		}
	}

	return parents, rows.Err()
}
