package account_test

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/profile/profile_search"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestProfileListOrderAndPaging(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ps *profile_search.Querier,
		raw *sqlx.DB,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			joined := []account.AccountID{}
			for _, s := range []account.Account{
				seed.Account_002_Frigg,
				seed.Account_003_Baldur,
				seed.Account_004_Loki,
			} {
				_, acc := e2e.WithAccount(root, aw, s)
				joined = append(joined, acc.ID)
			}

			list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{}, session)
			tests.Ok(t, err, list)

			positions := map[string]int{}
			for i, p := range list.JSON200.Profiles {
				positions[p.Id] = i
			}

			for i := 1; i < len(joined); i++ {
				previous, ok := positions[joined[i-1].String()]
				r.True(ok, "account %s missing from the profile list", joined[i-1])
				current, ok := positions[joined[i].String()]
				r.True(ok, "account %s missing from the profile list", joined[i])

				a.Less(current, previous, "the newest member comes first")
			}

			// one timestamp for every row, so only the id can separate them
			shared := time.Now().Add(-time.Hour).UTC()
			for _, id := range joined {
				_, err := raw.ExecContext(ctx, raw.Rebind("update accounts set created_at = ? where id = ?"), shared, id.String())
				r.NoError(err)
			}

			seen := map[account.AccountID]int{}

			for page := uint(1); page <= 10; page++ {
				result, err := ps.Search(root, pagination.NewPageParams(page, 2))
				r.NoError(err)
				r.Equal(int(page), result.CurrentPage)

				for _, p := range result.Items {
					seen[p.ID]++
				}

				if _, ok := result.NextPage.Get(); !ok {
					break
				}
			}

			for _, id := range joined {
				a.Equal(1, seen[id], "profile %s should appear on exactly one page", id)
			}
		}))
	}))
}
