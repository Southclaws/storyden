package role_repo

import (
	"context"
	"math"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_badge"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

const (
	// role data doesn't change often, and we don't want it to expire.
	cacheTTL                   = math.MaxInt64
	roleCachePrefix            = "role:by-id:"
	roleCustomOrderingCacheKey = "role:custom-ordering"
	accountRoleCachePrefix     = "role:account-assignments:"
)

type Repository struct {
	db    *ent.Client
	store cache.Store
}

func New(db *ent.Client, store cache.Store) *Repository {
	return &Repository{
		db:    db,
		store: store,
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			New,
			fx.Annotate(
				func(c *Repository) *Repository { return c },
				fx.As(new(Writer)),
			),
			fx.Annotate(
				func(c *Repository) *Repository { return c },
				fx.As(new(role_assign.Assign)),
			),
			fx.Annotate(
				func(c *Repository) *Repository { return c },
				fx.As(new(role_badge.Badge)),
			),
		),
		fx.Invoke(bindLifecycle),
	)
}

func bindLifecycle(lc fx.Lifecycle, repo *Repository) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		if err := repo.syncAll(ctx); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		return nil
	}))
}
