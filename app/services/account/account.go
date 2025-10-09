package account

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account/account_manage"
	"github.com/Southclaws/storyden/app/services/account/account_update"
	"github.com/Southclaws/storyden/app/services/account/profile_semdex"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(account_manage.New),
		fx.Provide(account_update.New),
		profile_semdex.Build(),
	)
}
