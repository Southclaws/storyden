package webauthn

import (
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/pkg/errors"
)

func New() (*webauthn.WebAuthn, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          "storyden.org",
		RPDisplayName: "Storyden",
		RPOrigin:      "localhost",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to init webauthn")
	}

	return wa, nil
}
