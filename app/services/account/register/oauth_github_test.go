package register_test

import (
	"context"
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestOAuthGitHubDuplicateAuthMethod(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		registrar *register.Registrar,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			service := authentication.ServiceOAuthGitHub
			token := "test-token"

			t.Run("duplicate_via_email_returns_same_account", func(t *testing.T) {
				authName := "GitHub (@emailuser)"
				identifierEmail := "email123"
				emailAddr := mail.Address{Address: "unique@example.com"}

				acc1, err := registrar.GetOrCreateViaEmail(root, service, authName, identifierEmail, token, "handleemail", "Email User", emailAddr)
				r.NoError(err)
				r.NotNil(acc1)

				acc2, err := registrar.GetOrCreateViaEmail(root, service, authName, identifierEmail, token, "handleemail", "Email User", emailAddr)
				r.NoError(err)
				r.NotNil(acc2)
				a.Equal(acc1.ID, acc2.ID)
			})

			t.Run("duplicate_via_handle", func(t *testing.T) {
				authName2 := "GitHub (@handleuser)"
				identifier2 := "54321"
				handle2 := "uniquehandle"

				acc1, err := registrar.GetOrCreateViaHandle(root, service, authName2, identifier2, token, handle2, "User Two")
				r.NoError(err)
				r.NotNil(acc1)

				acc2, err := registrar.GetOrCreateViaHandle(root, service, authName2, identifier2, token, handle2, "User Two")
				r.NoError(err)
				r.NotNil(acc2)
				a.Equal(acc1.ID, acc2.ID)
			})

			t.Run("duplicate_auth_method_via_random_handle", func(t *testing.T) {
				identifier3 := "99999"

				acc1, err := registrar.CreateWithRandomHandle(root, service, "GitHub (@randomuser)", identifier3, token, "Random User")
				r.NoError(err)
				r.NotNil(acc1)

				acc2, err := registrar.CreateWithRandomHandle(root, service, "GitHub (@randomuser)", identifier3, token, "Random User")
				r.Error(err)
				r.Nil(acc2)
				a.Contains(err.Error(), "authentication method already linked to another account")
			})

			t.Run("duplicate_auth_method_via_handle", func(t *testing.T) {
				identifier4 := "88888"
				handle4 := "specifichandle"

				acc1, err := registrar.CreateWithHandle(root, service, "GitHub (@specific)", identifier4, token, "Specific User", handle4)
				r.NoError(err)
				r.NotNil(acc1)

				acc2, err := registrar.CreateWithHandle(root, service, "GitHub (@specific)", identifier4, token, "Specific User", "differenthandle")
				r.Error(err)
				r.Nil(acc2)
				a.Contains(err.Error(), "authentication method already linked to another account")
			})
		}))
	}))
}
