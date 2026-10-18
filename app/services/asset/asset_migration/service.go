package asset_migration

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_repo"
	migration "github.com/Southclaws/storyden/app/resources/asset/asset_migration"
	jobs "github.com/Southclaws/storyden/app/resources/background_job"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/background_job"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Migrator struct {
	jobs     *jobs.Repository
	accounts *account_repo.Repository
	settings *settings.SettingsRepository
	objects  object.Storer
	runner   *background_job.Runner
	logger   *slog.Logger
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(func(*Migrator) {}),
	)
}

func New(
	lc fx.Lifecycle,
	ctx context.Context,
	jobs *jobs.Repository,
	accounts *account_repo.Repository,
	settings *settings.SettingsRepository,
	objects object.Storer,
	runner *background_job.Runner,
	logger *slog.Logger,
) (*Migrator, error) {
	m := &Migrator{
		jobs:     jobs,
		accounts: accounts,
		settings: settings,
		objects:  objects,
		runner:   runner,
		logger:   logger,
	}

	err := runner.Register(migration.Kind, 1, background_job.Handler[migration.Payload]{
		Run: m.run,
	})
	if err != nil {
		return nil, err
	}

	lc.Append(fx.StartHook(func() error {
		return m.scan(ctx)
	}))

	return m, nil
}
