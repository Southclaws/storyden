package account_repo

import (
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

type Repository struct {
	db           *ent.Client
	roleRepo     *role_repo.Repository
	roleHydrator *role_hydrate.Hydrator
	accountCache *accountCache
}

func New(db *ent.Client, roleRepo *role_repo.Repository, roleHydrator *role_hydrate.Hydrator, store cache.Store) *Repository {
	return &Repository{
		db:           db,
		roleRepo:     roleRepo,
		roleHydrator: roleHydrator,
		accountCache: newAccountCache(store),
	}
}
