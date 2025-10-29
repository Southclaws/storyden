package collection_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/collection/collection_querier"
	"github.com/Southclaws/storyden/app/resources/collection/collection_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/services/collection/collection_auth"
)

type Manager struct {
	colQuerier *collection_querier.Querier
	colWriter  *collection_writer.Writer
}

func New(
	colQuerier *collection_querier.Querier,
	colWriter *collection_writer.Writer,
) *Manager {
	return &Manager{
		colQuerier: colQuerier,
		colWriter:  colWriter,
	}
}

type Partial struct {
	Name        opt.Optional[string]
	Slug        opt.Optional[string]
	Description opt.Optional[string]
}

func (s *Manager) Create(ctx context.Context, accID account.AccountID, name string, partial Partial) (*collection.CollectionWithItems, error) {
	opts := []collection_writer.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, collection_writer.WithName(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, collection_writer.WithDescription(v)) })

	slug := partial.Slug.Or(mark.Slugify(name))

	col, err := s.colWriter.Create(ctx, accID, name, slug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *Manager) Update(ctx context.Context, qk collection.QueryKey, partial Partial) (*collection.CollectionWithItems, error) {
	if err := s.authoriseDirectUpdate(ctx, qk); err != nil {
		return nil, err
	}

	opts := []collection_writer.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, collection_writer.WithName(v)) })
	partial.Slug.Call(func(v string) { opts = append(opts, collection_writer.WithSlug(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, collection_writer.WithDescription(v)) })

	col, err := s.colWriter.Update(ctx, qk, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *Manager) Delete(ctx context.Context, qk collection.QueryKey) error {
	if err := s.authoriseDirectUpdate(ctx, qk); err != nil {
		return err
	}

	err := s.colWriter.Delete(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (m *Manager) authoriseDirectUpdate(ctx context.Context, qk collection.QueryKey) error {
	col, err := m.colQuerier.Probe(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionMutationPermissions(ctx, *col)
}
