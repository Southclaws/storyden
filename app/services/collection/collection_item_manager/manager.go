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
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/collection/collection_auth"
)

type Manager struct {
	colQuerier *collection_querier.Querier
	repo       *collection_item.Repository
}

func New(
	colQuerier *collection_querier.Querier,
	repo *collection_item.Repository,
) *Manager {
	return &Manager{
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
	acc, err := session.GetAccount(ctx)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	citem, err := m.repo.ProbeItem(ctx, qk, iid)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	mt, err := collection_auth.CheckCollectionItemMutationPermissions(ctx, acc, *citem)
	if err != nil {
		return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
	}

	if mt == collection.MembershipTypeSubmissionReview {
		if err := session.Authorise(ctx, nil, rbac.PermissionCollectionSubmit); err != nil {
			return collection.MembershipType{}, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return mt, nil
}

func (m *Manager) authoriseDirectUpdate(ctx context.Context, qk collection.QueryKey) error {
	col, err := m.colQuerier.Probe(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return collection_auth.CheckCollectionMutationPermissions(ctx, *col)
}
