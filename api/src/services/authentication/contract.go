package authentication

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/user"
)

type Contract interface {
	Encode(w http.ResponseWriter, user user.User) error
	Decode(r *http.Request) (*user.User, error)
}

func Build() fx.Option {
	return fx.Provide(NewCookieAuth)
}
