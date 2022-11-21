package securecookie

import (
	"github.com/gorilla/securecookie"

	"github.com/Southclaws/storyden/internal/config"
)

func New(cfg config.Config) *securecookie.SecureCookie {
	// TODO: Generate these on the fly if they do not already exist and warn
	// the administrator that if they aren't constant across server restarts,
	// cookies will break. Print out the keys on first run for convenience.
	return securecookie.New(cfg.HashKey, cfg.BlockKey)
}
