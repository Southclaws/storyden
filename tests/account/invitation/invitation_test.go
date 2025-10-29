package invitation_test

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestAccountInvitations(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		accountQuery *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			inviterCtx, inviter := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			inviterSession := sh.WithSession(inviterCtx)

			message := "Join me on Storyden!"
			invResponse, err := cl.InvitationCreateWithResponse(root, openapi.InvitationInitialProps{
				Message: &message,
			}, inviterSession)
			tests.Ok(t, err, invResponse)
			a.Equal(inviter.ID.String(), invResponse.JSON200.Creator.Id)
			r.NotNil(invResponse.JSON200.Message)
			a.Equal(message, *invResponse.JSON200.Message)

			invitationID := invResponse.JSON200.Id

			t.Run("accept_invite_with_password", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				inviteeHandle := "invitee-" + xid.New().String()
				inviteeResponse, err := cl.AuthPasswordSignupWithResponse(root, &openapi.AuthPasswordSignupParams{
					InvitationId: &invitationID,
				}, openapi.AuthPair{Identifier: inviteeHandle, Token: "password"})
				tests.Ok(t, err, inviteeResponse)
				id := account.AccountID(utils.Must(xid.FromString(inviteeResponse.JSON200.Id)))

				invitee, err := cl.AccountGetWithResponse(root, sh.WithSession(e2e.WithAccountID(root, id)))
				tests.Ok(t, err, invitee)
				r.NotNil(invitee.JSON200.InvitedBy)
				a.Equal(inviter.ID.String(), invitee.JSON200.InvitedBy.Id)

				public, err := cl.ProfileGetWithResponse(root, invitee.JSON200.Handle)
				tests.Ok(t, err, public)
				r.NotNil(public.JSON200.InvitedBy)
				a.Equal(inviter.ID.String(), public.JSON200.InvitedBy.Id)
			})

			t.Run("accept_invite_with_email_password", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				inviteeHandle := "invitee-" + xid.New().String()
				inviteeResponse, err := cl.AuthEmailPasswordSignupWithResponse(root, &openapi.AuthEmailPasswordSignupParams{
					InvitationId: &invitationID,
				}, openapi.AuthEmailPassword{Handle: &inviteeHandle, Email: xid.New().String() + "sc@storyden.org", Password: "password"})
				tests.Ok(t, err, inviteeResponse)
				id := account.AccountID(utils.Must(xid.FromString(inviteeResponse.JSON200.Id)))

				invitee, err := cl.AccountGetWithResponse(root, sh.WithSession(e2e.WithAccountID(root, id)))
				tests.Ok(t, err, invitee)
				r.NotNil(invitee.JSON200.InvitedBy)
				a.Equal(inviter.ID.String(), invitee.JSON200.InvitedBy.Id)

				public, err := cl.ProfileGetWithResponse(root, invitee.JSON200.Handle)
				tests.Ok(t, err, public)
				r.NotNil(public.JSON200.InvitedBy)
				a.Equal(inviter.ID.String(), public.JSON200.InvitedBy.Id)
			})
		}))
	}))
}
