package account_role

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_badge"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Manager struct {
	assign       role_assign.Assign
	badge        role_badge.Badge
	profileCache *profile_cache.Cache
	bus          *pubsub.Bus
}

func New(assign role_assign.Assign, badge role_badge.Badge, profileCache *profile_cache.Cache, bus *pubsub.Bus) *Manager {
	return &Manager{
		assign:       assign,
		badge:        badge,
		profileCache: profileCache,
		bus:          bus,
	}
}

func (m *Manager) UpdateRoles(ctx context.Context, accountID xid.ID, roles ...role_assign.Mutation) error {
	if err := m.profileCache.Invalidate(ctx, accountID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.assign.UpdateRoles(ctx, accountID, roles...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &message.EventAccountUpdated{
		ID: account.AccountID(accountID),
	})

	return nil
}

func (m *Manager) UpdateBadge(ctx context.Context, accountID xid.ID, roleID role.RoleID, badge bool) error {
	if err := m.profileCache.Invalidate(ctx, accountID); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := m.badge.UpdateBadge(ctx, accountID, roleID, badge); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &message.EventAccountUpdated{
		ID: account.AccountID(accountID),
	})

	return nil
}
