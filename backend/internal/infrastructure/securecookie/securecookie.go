package securecookie

import (
	"github.com/gorilla/securecookie"

	"github.com/Southclaws/storyden/backend/internal/config"
)

func New(cfg config.Config) *securecookie.SecureCookie {
	return securecookie.New(cfg.HashKey, cfg.BlockKey)
}
