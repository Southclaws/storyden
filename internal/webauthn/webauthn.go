package webauthn

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/Southclaws/storyden/internal/config"
)

func New(cfg config.Config) (*webauthn.WebAuthn, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Storyden",
		RPID:          cfg.CookieDomain,
		RPOrigin:      cfg.CookieDomain,
	})
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to init webauthn"))
	}

	return wa, nil
}
