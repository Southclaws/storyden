package securecookie_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/securecookie"
)

func TestSession(t *testing.T) {
	defer integration.Test(t, &config.Config{
		SessionKey: "e6386bdf71c523e4f313a244e64d560db172731674861e735a1eae90ecef98f3df6cb2a385f6d96cc9f9746ca83b8839868c",
	},
		fx.Invoke(func(
			e *securecookie.Session,
		) {
			a := assert.New(t)

			s := e.Encrypt("southclaws")
			a.NotEqual("southclaws", s)

			out, ok := e.Decrypt(s)
			a.True(ok)
			a.Equal("southclaws", out)

			s1 := e.Encrypt("southclaws")
			s2 := e.Encrypt("southclaws")
			s3 := e.Encrypt("southclaws")

			// assert nonce is suitably random.
			a.NotEqual(s1, s2)
			a.NotEqual(s2, s3)
			a.NotEqual(s3, s1)
		}),
	)
}
