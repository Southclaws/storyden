package collection_item_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item"
	"github.com/Southclaws/storyden/app/resources/collection/collection_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/collection/collection_auth"
)

type Manager struct {
	session    *session.Provider
	colQuerier *collection_querier.Querier
	repo       *collection_item.Repository
}

func New(
	session *session.Provider,
	colQuerier *collection_querier.Querier,
	repo *collection_item.Repository,
) *Manager {
	return &Manager{
		session:    session,
		colQuerier: colQuerier,
		repo:       repo,
	}
}

func (m *Manager) PostAdd(ctx context.Context, qk collection.QueryKey, pid post.ID) (*collection.CollectionWithItems, error) {
	mt, err := m.authoriseSubmission(ctx, qk, xid.ID(pid))
	if err != nil {
		return nil, err
	}

	col, err := m.repo.UpdateItems(ctx, qk, collection_item.WithPost(pid, mt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) PostRemove(ctx context.Context, qk collection.QueryKey, pid post.ID) (*collection.CollectionWithItems, error) {
	col, err := m.repo.UpdateItems(ctx, qk, collection_item.WithPostRemove(pid))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) NodeAdd(ctx context.Context, qk collection.QueryKey, id library.NodeID) (*collection.CollectionWithItems, error) {
	mt, err := m.authoriseSubmission(ctx, qk, xid.ID(id))
	if err != nil {
		return nil, err
	}

	col, err := m.repo.UpdateItems(ctx, qk, collection_item.WithNode(id, mt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) NodeRemove(ctx context.Context, qk collection.QueryKey, id library.NodeID) (*collection.CollectionWithItems, error) {
	if err := m.authoriseDirectUpdate(ctx, qk); err != nil {
		return nil, err
	}

	col, err := m.repo.UpdateItems(ctx, qk, collection_item.WithNodeRemove(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) authoriseSubmission(ctx context.Context, qk collection.QueryKey, iid xid.ID) (collection.MembershipType, error) {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	citem, err := m.repo.ProbeItem(ctx, qk, iid)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionItemMutationPermissions(ctx, *acc, *citem)
}

func (m *Manager) authoriseDirectUpdate(ctx context.Context, qk collection.QueryKey) error {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	col, err := m.colQuerier.Probe(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionMutationPermissions(ctx, *acc, *col)
}
