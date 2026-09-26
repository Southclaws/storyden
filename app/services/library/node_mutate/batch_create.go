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

func (s *Manager) CreateMany(ctx context.Context, items []BatchCreate) ([]BatchResult, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageLibrary); err != nil {
		return nil, err
	}

	owner, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	fields, err := prepareCreateFields(items)
	if err != nil {
		return nil, err
	}

	parents, err := s.resolveCreateParents(ctx, items, fields)
	if err != nil {
		return nil, err
	}

	options := s.batchOptions(ctx, fields, make([]node_querier.MutationState, len(items)), parents)
	writes := make([]node_writer.CreateInput, len(items))
	results := make([]BatchResult, len(items))
	events := make([]any, 0, len(items)*2)
	for i, item := range items {
		field := fields[i]
		slug := field.Slug.OrZero()
		writes[i] = node_writer.CreateInput{
			ID:      item.ID,
			Name:    field.Name.OrZero(),
			Slug:    slug,
			Options: options[i],
			Tags:    field.Tags,
		}
		results[i] = BatchResult{Mark: library.NewMark(xid.ID(item.ID), slug.String()), Name: field.Name.OrZero()}
		events = append(events, &rpc.EventNodeCreated{ID: item.ID, Slug: slug.String()})
		if field.Visibility.OrZero() == visibility.VisibilityPublished {
			events = append(events, &rpc.EventNodePublished{ID: item.ID, Slug: slug.String()})
		}
	}

	if err := s.cache.Invalidate(ctx); err != nil {
		return nil, err
	}

	if err := s.nodeWriter.CreateMany(ctx, owner, writes); err != nil {
		return nil, fmt.Errorf("create page batch: %w", err)
	}
	s.cache.InvalidateAfterWrite(ctx)

	s.bus.PublishMany(ctx, events...)

	return results, nil
}

func prepareCreateFields(items []BatchCreate) ([]BatchFields, error) {
	fields := make([]BatchFields, len(items))
	allocated := make(map[library.NodeID]bool, len(items))
	slugs := make(map[string]library.NodeID, len(items))
	for i, item := range items {
		if item.ID == (library.NodeID{}) || allocated[item.ID] {
			return nil, fmt.Errorf("item %d: ID must be nonzero and unique", i+1)
		}

		field := item.Fields
		if !field.Name.Ok() {
			return nil, fmt.Errorf("item %d: name is required", i+1)
		}

		field.Slug = opt.New(field.Slug.Or(mark.NewSlugFromName(field.Name.OrZero())))
		field.Visibility = opt.New(field.Visibility.Or(visibility.VisibilityDraft))
		slug := field.Slug.OrZero().String()
		if _, exists := slugs[slug]; exists {
			return nil, fmt.Errorf("item %d: duplicate slug %q", i+1, slug)
		}

		allocated[item.ID] = true
		slugs[slug] = item.ID
		fields[i] = field
	}

	if err := validateBatchFields(fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func (s *Manager) resolveCreateParents(ctx context.Context, items []BatchCreate, fields []BatchFields) ([]opt.Optional[library.NodeID], error) {
	allocated := make(map[library.NodeID]bool, len(items))
	slugs := make(map[string]bool, len(items))
	for i, item := range items {
		allocated[item.ID] = true
		slugs[fields[i].Slug.OrZero().String()] = true
	}

	keys := make([]library.QueryKey, 0, len(items)*3)
	requestedSlugs := make([]mark.Slug, len(items))
	for i, item := range items {
		keys = append(keys, library.NewID(xid.ID(item.ID)))
		requestedSlugs[i] = fields[i].Slug.OrZero()
		if parent, ok := fields[i].Parent.Get(); ok {
			id, hasID := parent.ID().Get()
			if !hasID || !allocated[library.NodeID(id)] {
				keys = append(keys, parent)
			}
		}
	}

	states, err := s.nodeQuerier.MutationStates(ctx, keys, requestedSlugs...)
	if err != nil {
		return nil, err
	}

	for id, state := range states {
		if allocated[id] {
			return nil, fmt.Errorf("page ID %s already exists", id)
		}

		if _, exists := slugs[state.Mark.Slug()]; exists {
			return nil, fmt.Errorf("slug %q already exists", state.Mark.Slug())
		}
	}

	index := indexMutationStates(states)
	for _, item := range items {
		states[item.ID] = node_querier.MutationState{}
	}

	parents := make([]opt.Optional[library.NodeID], len(items))
	graph := make(map[library.NodeID]library.NodeID)
	for i, field := range fields {
		if key, ok := field.Parent.Get(); ok {
			parent, err := resolveBatchKey(key, states, index)
			if err != nil {
				return nil, fmt.Errorf("item %d parent: %w", i+1, err)
			}

			parents[i] = opt.New(parent)
			graph[items[i].ID] = parent
		}
	}

	if err := validateBatchHierarchy(graph); err != nil {
		return nil, err
	}

	return parents, nil
}
