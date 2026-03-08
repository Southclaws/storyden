package account_role_assign

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Manager struct {
	assign       *role_assign.Assignment
	profileCache *profile_cache.Cache
	bus          *pubsub.Bus
}

func New(assign *role_assign.Assignment, profileCache *profile_cache.Cache, bus *pubsub.Bus) *Manager {
	return &Manager{
		assign:       assign,
		profileCache: profileCache,
		bus:          bus,
	}
}

func (m *Manager) UpdateRoles(ctx context.Context, accountID account_ref.ID, roles ...role_assign.Mutation) error {
	if err := m.profileCache.Invalidate(ctx, xid.ID(accountID)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.assign.UpdateRoles(ctx, accountID, roles...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountUpdated{ID: accountID})

	return nil
}

func (m *Manager) SetBadge(ctx context.Context, accountID account_ref.ID, roleID role.RoleID, badge bool) error {
	if err := m.profileCache.Invalidate(ctx, xid.ID(accountID)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.assign.SetBadge(ctx, accountID, roleID, badge); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountUpdated{ID: accountID})

	return nil
}
