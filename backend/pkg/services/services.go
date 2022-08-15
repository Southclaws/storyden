package services

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/services/account"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication"
	"github.com/Southclaws/storyden/backend/pkg/services/rbac"
)

func Build() fx.Option {
	return fx.Options(
		account.Build(),
		authentication.Build(),
		rbac.Build(),
	)
}
