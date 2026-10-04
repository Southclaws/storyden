package oauth_test

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/oauthregistrationapproval"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthRegistrationManagement(t *testing.T) {
	t.Parallel()
	integration.Test(t, oauthConfig(t), e2e.Setup(), withOAuthRegistration(t, "approval"), fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			memberCtx, _ := e2e.WithAccount(ctx, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)
			now := time.Now()

			t.Run("tokens_paginate_without_leaking_credentials", func(t *testing.T) {
				expected := make([]string, 100)
				for n := range 100 {
					token, err := db.OAuthDynamicRegistrationAccessTokens.Create().SetLabel(fmt.Sprintf("Token %03d", n)).SetTokenHash("secret-hash-" + xid.New().String()).SetCreatorAccountID(xid.ID(admin.ID)).SetExpiresAt(now.Add(time.Hour)).SetMaxRegistrations(1).SetCreatedAt(now.Add(time.Duration(n) * time.Millisecond)).Save(ctx)
					require.NoError(t, err)
					expected[n] = token.ID.String()
				}
				actual := []string{}
				for pageNumber := 1; ; pageNumber++ {
					page := strconv.Itoa(pageNumber)
					response := tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenListWithResponse(ctx, &openapi.AdminOAuthDynamicRegistrationAccessTokenListParams{Page: &page}, adminSession))(t, http.StatusOK)
					require.Equal(t, 50, response.JSON200.PageSize)
					require.LessOrEqual(t, len(response.JSON200.Iats), 50)
					require.NotContains(t, string(response.Body), "secret-hash")
					for _, token := range response.JSON200.Iats {
						actual = append(actual, string(token.Id))
					}

					if pageNumber >= response.JSON200.TotalPages {
						break
					}
				}
				for _, id := range expected {
					require.Contains(t, actual, id)
				}
			})

			create := func(count int, createdAt, expiresAt time.Time) []*ent.OAuthRegistrationApproval {
				rows := make([]*ent.OAuthRegistrationApproval, 0, count)
				for range count {
					id := xid.New().String()
					row, err := db.OAuthRegistrationApproval.Create().SetRegistrationCodeHash(id).SetVerificationCodeHash(id).SetVerificationCodeDisplay(id).SetMetadata([]byte(`{"client_name":"test-bot"}`)).SetCreatedAt(createdAt).SetExpiresAt(expiresAt).SetNextPollAt(now).Save(ctx)
					require.NoError(t, err)
					rows = append(rows, row)
				}
				return rows
			}
			t.Run("approve_all_and_deny_all_use_pending_snapshot_across_pages", func(t *testing.T) {
				team := create(25, now, now.Add(time.Hour))
				queue := tests.AssertRequest(cl.AdminOAuthRegistrationApprovalListWithResponse(ctx, nil, adminSession))(t, http.StatusOK).JSON200
				approval := openapi.AdminOAuthRegistrationApprovalBulkSubmitJSONRequestBody{Approved: true, CreatedBefore: queue.SnapshotAt.UTC()}
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalBulkSubmitWithResponse(ctx, approval))(t, http.StatusUnauthorized)
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalBulkSubmitWithResponse(ctx, approval, memberSession))(t, http.StatusForbidden)
				accountsBefore, err := db.Account.Query().Count(ctx)
				require.NoError(t, err)
				approved := tests.AssertRequest(cl.AdminOAuthRegistrationApprovalBulkSubmitWithResponse(ctx, approval, adminSession))(t, http.StatusOK)
				require.Equal(t, 25, approved.JSON200.Updated)
				for _, row := range team {
					stored, err := db.OAuthRegistrationApproval.Get(ctx, row.ID)
					require.NoError(t, err)
					require.NotNil(t, stored.ApprovedAt)
					require.Equal(t, xid.ID(admin.ID), *stored.ApprovedByAccountID)
				}
				accountsAfter, err := db.Account.Query().Count(ctx)
				require.NoError(t, err)
				require.Equal(t, accountsBefore, accountsAfter)

				excluded := append(create(1, now, now.Add(-time.Second)), create(3, now, now.Add(time.Hour))...)
				require.NoError(t, db.OAuthRegistrationApproval.UpdateOneID(excluded[1].ID).SetCancelledAt(now).Exec(ctx))
				require.NoError(t, db.OAuthRegistrationApproval.UpdateOneID(excluded[2].ID).SetConsumedAt(now).Exec(ctx))
				require.NoError(t, db.OAuthRegistrationApproval.UpdateOneID(excluded[3].ID).SetDeniedAt(now).Exec(ctx))
				spam := create(1000, now, now.Add(time.Hour))
				queue = tests.AssertRequest(cl.AdminOAuthRegistrationApprovalListWithResponse(ctx, nil, adminSession))(t, http.StatusOK).JSON200
				require.Equal(t, 20, queue.TotalPages)
				require.Len(t, queue.Registrations, 50)
				late := create(1, queue.SnapshotAt.Add(time.Second), now.Add(time.Hour))[0]
				denial := openapi.AdminOAuthRegistrationApprovalBulkSubmitJSONRequestBody{Approved: false, CreatedBefore: queue.SnapshotAt}
				denied := tests.AssertRequest(cl.AdminOAuthRegistrationApprovalBulkSubmitWithResponse(ctx, denial, adminSession))(t, http.StatusOK)
				require.Equal(t, 1000, denied.JSON200.Updated)
				ids := make([]xid.ID, len(spam))
				for n, row := range spam {
					ids[n] = row.ID
				}
				count, err := db.OAuthRegistrationApproval.Query().Where(oauthregistrationapproval.IDIn(ids...), oauthregistrationapproval.DeniedAtNotNil()).Count(ctx)
				require.NoError(t, err)
				require.Equal(t, 1000, count)
				for _, row := range append(team, excluded[:3]...) {
					stored, err := db.OAuthRegistrationApproval.Get(ctx, row.ID)
					require.NoError(t, err)
					require.Nil(t, stored.DeniedAt)
				}
				unchanged, err := db.OAuthRegistrationApproval.Get(ctx, late.ID)
				require.NoError(t, err)
				require.Nil(t, unchanged.DeniedAt)
				require.Nil(t, unchanged.ApprovedAt)
				repeat := tests.AssertRequest(cl.AdminOAuthRegistrationApprovalBulkSubmitWithResponse(ctx, denial, adminSession))(t, http.StatusOK)
				require.Zero(t, repeat.JSON200.Updated)
				future := denial
				future.CreatedBefore = time.Now().Add(time.Hour)
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalBulkSubmitWithResponse(ctx, future, adminSession))(t, http.StatusBadRequest)
			})
		}))
	}))
}
