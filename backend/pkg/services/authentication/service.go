package authentication

import (
	"net/http"

	"github.com/Southclaws/storyden/backend/pkg/resources/user"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider"
	"go.uber.org/fx"
)

type Service interface {
	DecodeSession(r *http.Request) (*user.User, bool)
}

func Build() fx.Option {
	return fx.Options(
		provider.Build(),
	)
}
