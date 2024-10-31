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
	"github.com/Southclaws/storyden/app/services/account/session"
	"github.com/Southclaws/storyden/app/services/collection/collection_auth"
)

type Manager struct {
	session    session.SessionProvider
	colQuerier *collection_querier.Querier
	colWriter  *collection_writer.Writer
}

func New(
	session session.SessionProvider,
	colQuerier *collection_querier.Querier,
	colWriter *collection_writer.Writer,
) *Manager {
	return &Manager{
		session:    session,
		colQuerier: colQuerier,
		colWriter:  colWriter,
	}
}

type Partial struct {
	Name        opt.Optional[string]
	Description opt.Optional[string]
}

func (s *Manager) Create(ctx context.Context, accID account.AccountID, name string, partial Partial) (*collection.CollectionWithItems, error) {
	opts := []collection_writer.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, collection_writer.WithName(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, collection_writer.WithDescription(v)) })

	col, err := s.colWriter.Create(ctx, accID, name, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *Manager) Update(ctx context.Context, cid collection.CollectionID, partial Partial) (*collection.CollectionWithItems, error) {
	if err := s.authoriseDirectUpdate(ctx, cid); err != nil {
		return nil, err
	}

	opts := []collection_writer.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, collection_writer.WithName(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, collection_writer.WithDescription(v)) })

	col, err := s.colWriter.Update(ctx, cid, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *Manager) Delete(ctx context.Context, cid collection.CollectionID) error {
	if err := s.authoriseDirectUpdate(ctx, cid); err != nil {
		return err
	}

	err := s.colWriter.Delete(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (m *Manager) authoriseDirectUpdate(ctx context.Context, cid collection.CollectionID) error {
	col, err := m.colQuerier.Probe(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.session.Account(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionMutationPermissions(ctx, *acc, *col)
}
