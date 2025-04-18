package token

import (
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			func(db *ent.Client, cs cache.Store) Repository {
				repo := New(db)

				return NewCachedRepository(repo, cs)
			},
		),
	)
}
