package account_repo

import (
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Repository struct {
	ins          spanner.Instrumentation
	db           *ent.Client
	roleRepo     *role_repo.Repository
	roleHydrator *role_hydrate.Hydrator
	accountCache *accountCache
}

func New(ins spanner.Builder, db *ent.Client, roleRepo *role_repo.Repository, roleHydrator *role_hydrate.Hydrator, store cache.Store) *Repository {
	return &Repository{
		ins:          ins.Build(),
		db:           db,
		roleRepo:     roleRepo,
		roleHydrator: roleHydrator,
		accountCache: newAccountCache(store),
	}
}
