package webauthn

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/duo-labs/webauthn/webauthn"
)

func New() (*webauthn.WebAuthn, error) {
	// TODO: Read this from config.
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          "storyden.org",
		RPDisplayName: "Storyden",
		RPOrigin:      "localhost",
	})
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to init webauthn"))
	}

	return wa, nil
}
