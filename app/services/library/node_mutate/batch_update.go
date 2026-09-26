package node_mutate

import (
	"context"
	"fmt"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (s *Manager) UpdateMany(ctx context.Context, items []BatchUpdate) ([]BatchResult, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageLibrary); err != nil {
		return nil, err
	}

	fields := make([]BatchFields, len(items))
	keys := make([]library.QueryKey, 0, len(items)*3)
	slugs := make([]mark.Slug, 0, len(items))
	for i, item := range items {
		fields[i] = item.Fields
		keys = append(keys, item.Key)
		item.Fields.Parent.Call(func(key library.QueryKey) { keys = append(keys, key) })
		if slug, ok := item.Fields.Slug.Get(); ok {
			slugs = append(slugs, slug)
		}
	}

	if err := validateBatchFields(fields); err != nil {
		return nil, err
	}

	states, err := s.nodeQuerier.MutationStates(ctx, keys, slugs...)
	if err != nil {
		return nil, err
	}

	current, parents, err := resolveUpdateTargets(items, states)
	if err != nil {
		return nil, err
	}

	ids := make([]library.NodeID, len(current))
	for i, state := range current {
		ids[i] = library.NodeID(state.Mark.ID())
	}

	if err := s.validateUpdateHierarchy(ctx, ids, parents); err != nil {
		return nil, err
	}

	options := s.batchOptions(ctx, fields, current, parents)
	writes := make([]node_writer.UpdateInput, len(items))
	results := make([]BatchResult, len(items))
	events := make([]any, 0, len(items)*2)
	for i, id := range ids {
		field, state := fields[i], current[i]
		slug := state.Mark.Slug()
		if value, ok := field.Slug.Get(); ok {
			slug = value.String()
		}

		if field.versioned() {
			options[i] = append(options[i], node_writer.WithCurrentVersionCleared())
		}

		writes[i] = node_writer.UpdateInput{
			ID:        id,
			UpdatedAt: state.UpdatedAt,
			Versioned: field.versioned(),
			Options:   options[i],
			Tags:      field.Tags,
		}
		results[i] = BatchResult{Mark: library.NewMark(xid.ID(id), slug), Name: field.Name.Or(state.Name)}
		events = append(events, &rpc.EventNodeUpdated{ID: id, Slug: slug})
		events = append(events, batchVisibilityEvents(id, slug, state.Visibility, field.Visibility.Or(state.Visibility))...)
	}

	if err := s.cache.InvalidateBeforeWrite(ctx); err != nil {
		return nil, err
	}

	err = s.nodeWriter.UpdateMany(ctx, writes)
	s.cache.InvalidateAfterWrite(ctx)
	if err != nil {
		return nil, fmt.Errorf("update page batch: %w", err)
	}

	s.bus.PublishMany(ctx, events...)

	return results, nil
}

func resolveUpdateTargets(items []BatchUpdate, states map[library.NodeID]node_querier.MutationState) ([]node_querier.MutationState, []opt.Optional[library.NodeID], error) {
	index := indexMutationStates(states)
	current := make([]node_querier.MutationState, len(items))
	parents := make([]opt.Optional[library.NodeID], len(items))
	seen := make(map[library.NodeID]bool, len(items))
	desiredSlugs := make(map[string]bool, len(items))
	for i, item := range items {
		id, err := resolveBatchKey(item.Key, states, index)
		if err != nil {
			return nil, nil, fmt.Errorf("item %d target: %w", i+1, err)
		}

		if seen[id] {
			return nil, nil, fmt.Errorf("multiple updates target page %s", id)
		}

		state := states[id]
		if state.HasDraft && item.Fields.versioned() {
			return nil, nil, fmt.Errorf("item %d: page %s has a working draft; apply or delete it before editing versioned fields", i+1, id)
		}

		slug := state.Mark.Slug()
		if value, ok := item.Fields.Slug.Get(); ok {
			slug = value.String()
		}

		if other, exists := index[slug]; (exists && other != id) || desiredSlugs[slug] {
			return nil, nil, fmt.Errorf("item %d: slug %q already exists", i+1, slug)
		}

		if key, ok := item.Fields.Parent.Get(); ok {
			parent, err := resolveBatchKey(key, states, index)
			if err != nil {
				return nil, nil, fmt.Errorf("item %d parent: %w", i+1, err)
			}

			parents[i] = opt.New(parent)
		}

		seen[id], desiredSlugs[slug] = true, true
		current[i] = state
	}

	return current, parents, nil
}

func (s *Manager) validateUpdateHierarchy(ctx context.Context, ids []library.NodeID, parents []opt.Optional[library.NodeID]) error {
	roots := make([]library.NodeID, 0, len(parents))
	for _, parent := range parents {
		if id, ok := parent.Get(); ok {
			roots = append(roots, id)
		}
	}

	if len(roots) == 0 {
		return nil
	}

	graph, err := s.nodeQuerier.Ancestors(ctx, roots)
	if err != nil {
		return err
	}

	for i, parent := range parents {
		if id, ok := parent.Get(); ok {
			graph[ids[i]] = id
		}
	}

	return validateBatchHierarchy(graph)
}

func batchVisibilityEvents(id library.NodeID, slug string, before, after visibility.Visibility) []any {
	if before == after {
		return nil
	}

	switch after {
	case visibility.VisibilityPublished:
		return []any{&rpc.EventNodePublished{ID: id, Slug: slug}}

	case visibility.VisibilityReview:
		return []any{&rpc.EventNodeSubmittedForReview{ID: id, Slug: slug}}

	default:
		if before == visibility.VisibilityPublished {
			return []any{&rpc.EventNodeUnpublished{ID: id, Slug: slug}}
		}
	}

	return nil
}
