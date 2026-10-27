package settings_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Manager struct {
	repo   *settings.SettingsRepository
	bus    *pubsub.Bus
	assets *asset_querier.Querier
}

func New(repo *settings.SettingsRepository, bus *pubsub.Bus, assets *asset_querier.Querier) *Manager {
	return &Manager{
		repo:   repo,
		bus:    bus,
		assets: assets,
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

	if err := s.Services.OrZero().Media.OrZero().Validate(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if err := m.validateBanner(ctx, s); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := m.repo.Set(ctx, s)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventSettingsUpdated{
		Settings: rpc.SerialiseSettings(*updated),
	})

	return updated, nil
}

func (m *Manager) validateBanner(ctx context.Context, s settings.Settings) error {
	id, ok := s.BannerAssetID.Get()
	if !ok {
		return nil
	}

	image, err := m.assets.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !asset.ImageEligible(image.MIME.String()) {
		return fault.Wrap(fault.New("banner must be a PNG, JPEG or WebP image"), ftag.With(ftag.InvalidArgument))
	}
	return nil
}
