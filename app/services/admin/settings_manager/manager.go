package settings_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Manager struct {
	repo *settings.SettingsRepository
	bus  *pubsub.Bus
}

func New(repo *settings.SettingsRepository, bus *pubsub.Bus) *Manager {
	return &Manager{
		repo: repo,
		bus:  bus,
	}
}

func (m *Manager) Get(ctx context.Context) (*settings.Settings, error) {
	s, err := m.repo.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s, nil
}

func (m *Manager) Set(ctx context.Context, s settings.Settings) (*settings.Settings, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageSettings, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := m.repo.Set(ctx, s)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &message.EventSettingsUpdated{
		Settings: updated,
	})

	return updated, nil
}
