package account_email_test

import (
	"context"
	"net/mail"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestManualVerifiedStatusSurvivesEmailVerification(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		aw *account_writer.Writer,
		aq *account_querier.Querier,
		er *email.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			ctx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)

			_, err := aw.Update(ctx, member.ID, account_writer.SetVerifiedStatus(account.VerifiedStatusManual))
			r.NoError(err)

			addr := mail.Address{Address: "loki-manual@storyden.org"}
			_, err = er.Add(ctx, member.ID, addr, "123456")
			r.NoError(err)

			err = er.Verify(ctx, member.ID, addr, "123456")
			r.NoError(err)

			acc, err := aq.GetByID(ctx, member.ID)
			r.NoError(err)
			r.Equal(account.VerifiedStatusManual, acc.VerifiedStatus, "manual admin override must survive email verification")
		}))
	}))
}
