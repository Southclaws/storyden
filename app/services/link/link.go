package link

import (
	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/services/link/hydrator"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/app/services/link/scrape_job"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			fetcher.New,
			hydrator.New,
			scrape.New,
		),
		scrape_job.Build(),
	)
}
