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
	"github.com/Southclaws/storyden/app/services/account/session"
	"github.com/Southclaws/storyden/app/services/collection/collection_auth"
)

type Manager struct {
	session    session.SessionProvider
	colQuerier *collection_querier.Querier
	repo       *collection_item.Repository
}

func New(
	session session.SessionProvider,
	colQuerier *collection_querier.Querier,
	repo *collection_item.Repository,
) *Manager {
	return &Manager{
		session:    session,
		colQuerier: colQuerier,
		repo:       repo,
	}
}

func (m *Manager) PostAdd(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.CollectionWithItems, error) {
	mt, err := m.authoriseSubmission(ctx, cid, xid.ID(pid))
	if err != nil {
		return nil, err
	}

	col, err := m.repo.UpdateItems(ctx, cid, collection_item.WithPost(pid, mt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) PostRemove(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.CollectionWithItems, error) {
	col, err := m.repo.UpdateItems(ctx, cid, collection_item.WithPostRemove(pid))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) NodeAdd(ctx context.Context, cid collection.CollectionID, id library.NodeID) (*collection.CollectionWithItems, error) {
	mt, err := m.authoriseSubmission(ctx, cid, xid.ID(id))
	if err != nil {
		return nil, err
	}

	col, err := m.repo.UpdateItems(ctx, cid, collection_item.WithNode(id, mt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) NodeRemove(ctx context.Context, cid collection.CollectionID, id library.NodeID) (*collection.CollectionWithItems, error) {
	if err := m.authoriseDirectUpdate(ctx, cid); err != nil {
		return nil, err
	}

	col, err := m.repo.UpdateItems(ctx, cid, collection_item.WithNodeRemove(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (m *Manager) authoriseSubmission(ctx context.Context, cid collection.CollectionID, iid xid.ID) (collection.MembershipType, error) {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	citem, err := m.repo.ProbeItem(ctx, cid, iid)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionItemMutationPermissions(ctx, *acc, *citem)
}

func (m *Manager) authoriseDirectUpdate(ctx context.Context, cid collection.CollectionID) error {
	acc, err := m.session.Account(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	col, err := m.colQuerier.Probe(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionMutationPermissions(ctx, *acc, *col)
}
