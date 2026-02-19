package account

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/account/account_manage"
	"github.com/Southclaws/storyden/app/services/account/account_role"
	"github.com/Southclaws/storyden/app/services/account/account_update"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(account_manage.New),
		fx.Provide(account_role.New),
		fx.Provide(account_update.New),
	)
}
