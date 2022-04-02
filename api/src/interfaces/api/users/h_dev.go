package users

import (
	"errors"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/Southclaws/storyden/api/src/infra/web"
)

type devParams struct {
	ID     string `qstring:"id"`
	Secret string `qstring:"secret"`
}

func (c *controller) dev(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("DEV_LOGIN_SECRET")
	if secret == "" {
		web.StatusNotAcceptable(w, errors.New("dev login not enabled"))
		return
	}

	var p devParams
	if !web.ParseQuery(w, r, &p) {
		web.StatusBadRequest(w, errors.New("invalid query"))
		return
	}

	if p.Secret != secret {
		web.StatusUnauthorized(w, errors.New("go away"))
		return
	}

	zap.L().Info("dev login", zap.Any("params", p))

	user, err := c.repo.GetUser(r.Context(), p.ID, false)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	zap.L().Info("dev login", zap.Any("user", user))

	c.auth.EncodeAuthCookie(w, *user)
	web.Write(w, user)
}
