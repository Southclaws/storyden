package account_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestPublicProfiles(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		ar *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			// Create 5 fresh accounts
			odin1 := newAccount(t, root, cl, ar, "odin")
			frigg := newAccount(t, root, cl, ar, "frigg")
			baldur := newAccount(t, root, cl, ar, "baldur")
			odin2 := newAccount(t, root, cl, ar, "odin2")
			thor := newAccount(t, root, cl, ar, "þórr")

			// Get them all, default params.

			list1, err := cl.ProfileListWithResponse(root, nil)
			r.NoError(err)
			r.NotNil(list1)
			r.Equal(http.StatusOK, list1.StatusCode())

			a.Equal(1, list1.JSON200.CurrentPage)
			a.GreaterOrEqual(len(list1.JSON200.Profiles), 5)

			// Get one specific name - search for "odin" should include both odins but exclude others

			list2, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
				Q: opt.New("odin").Ptr(),
			})
			r.NoError(err)
			r.NotNil(list2)
			r.Equal(http.StatusOK, list2.StatusCode())

			a.Equal(1, list2.JSON200.CurrentPage)
			a.Equal(50, list2.JSON200.PageSize)
			a.GreaterOrEqual(len(list2.JSON200.Profiles), 2)

			odinResults := findProfilesByNamePrefix(list2.JSON200.Profiles, "odin")
			a.GreaterOrEqual(len(odinResults), 2, "should find at least 2 odin profiles")

			_, foundOdin1 := findProfile(list2.JSON200.Profiles, odin1.Handle)
			_, foundOdin2 := findProfile(list2.JSON200.Profiles, odin2.Handle)
			_, foundFrigg := findProfile(list2.JSON200.Profiles, frigg.Handle)
			_, foundBaldur := findProfile(list2.JSON200.Profiles, baldur.Handle)
			_, foundThor := findProfile(list2.JSON200.Profiles, thor.Handle)

			a.True(foundOdin1, "odin1 should be in search results")
			a.True(foundOdin2, "odin2 should be in search results")
			a.False(foundFrigg, "frigg should not be in odin search results")
			a.False(foundBaldur, "baldur should not be in odin search results")
			a.False(foundThor, "thor should not be in odin search results")

			// Query an invalid page.

			list3, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
				Page: opt.New("2147483647").Ptr(),
			})
			r.NoError(err)
			r.NotNil(list3)
			r.Equal(http.StatusOK, list3.StatusCode())

			a.Equal(2147483647, list3.JSON200.CurrentPage)
			a.Nil(list3.JSON200.NextPage)
			a.Empty(list3.JSON200.Profiles)
		}))
	}))
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accountQuery *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			handle1 := "user-" + xid.New().String()
			acc1, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle1, "password"})
			tests.Ok(t, err, acc1)
			session1 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc1.JSON200.Id)))))

			handle2 := "user-" + xid.New().String()
			acc2, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{handle2, "password"})
			tests.Ok(t, err, acc2)
			session2 := sh.WithSession(e2e.WithAccountID(root, account.AccountID(utils.Must(xid.FromString(acc2.JSON200.Id)))))

			t.Run("update_profile", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// as guest
				get1, err := cl.ProfileGetWithResponse(root, handle1)
				tests.Ok(t, err, get1)
				r.Equal(handle1, get1.JSON200.Handle)
				r.Equal(handle1, get1.JSON200.Name)

				// as another user
				get2, err := cl.ProfileGetWithResponse(root, handle1, session2)
				tests.Ok(t, err, get2)
				r.Equal(handle1, get2.JSON200.Handle)
				r.Equal(handle1, get2.JSON200.Name)

				// update account profile
				newbio := "new bio"
				newhandle := "newhandle-" + xid.New().String()
				newname := "newname"
				newlinks := []openapi.ProfileExternalLink{
					{Text: "link1", Url: "https://example.com"},
				}
				newmeta := openapi.Metadata{
					"some": "data",
				}

				upd1, err := cl.AccountUpdateWithResponse(root, openapi.AccountUpdateJSONRequestBody{
					Bio:    &newbio,
					Handle: &newhandle,
					Name:   &newname,
					Links:  &newlinks,
					Meta:   &newmeta,
				}, session1)
				tests.Ok(t, err, upd1)

				a.Contains(upd1.JSON200.Bio, newbio)
				a.Equal(newhandle, upd1.JSON200.Handle)
				a.Equal(newname, upd1.JSON200.Name)
				r.Len(newlinks, 1)
				a.Equal(newlinks, upd1.JSON200.Links)
				a.Equal(newmeta, upd1.JSON200.Meta)

				// old handle should not work
				getold, err := cl.ProfileGetWithResponse(root, handle1)
				tests.Status(t, err, getold, http.StatusNotFound)

				// as guest
				getAfterUpdateAsGuest, err := cl.ProfileGetWithResponse(root, newhandle)
				tests.Ok(t, err, getAfterUpdateAsGuest)
				a.Contains(getAfterUpdateAsGuest.JSON200.Bio, newbio)
				a.Equal(newhandle, getAfterUpdateAsGuest.JSON200.Handle)
				a.Equal(newname, getAfterUpdateAsGuest.JSON200.Name)
				r.Len(newlinks, 1)
				a.Equal(newlinks, getAfterUpdateAsGuest.JSON200.Links)
				a.Equal(newmeta, getAfterUpdateAsGuest.JSON200.Meta)

				// as another user
				getAfterUpdateAsUser2, err := cl.ProfileGetWithResponse(root, newhandle, session2)
				tests.Ok(t, err, getAfterUpdateAsUser2)
				a.Contains(getAfterUpdateAsUser2.JSON200.Bio, newbio)
				a.Equal(newhandle, getAfterUpdateAsUser2.JSON200.Handle)
				a.Equal(newname, getAfterUpdateAsUser2.JSON200.Name)
				r.Len(newlinks, 1)
				a.Equal(newlinks, getAfterUpdateAsUser2.JSON200.Links)
				a.Equal(newmeta, getAfterUpdateAsUser2.JSON200.Meta)
			})
		}))
	}))
}

func TestProfileListFilters(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		ar *account_querier.Querier,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			prefix := "pft" + xid.New().String()[:3]
			acc1 := newAccount(t, root, cl, ar, prefix+"a")
			acc2 := newAccount(t, root, cl, ar, prefix+"b")
			acc3 := newAccount(t, root, cl, ar, prefix+"c")

			t.Run("sort_by_name_ascending", func(t *testing.T) {
				sortParam := "name"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Sort: &sortParam,
					Q:    &prefix,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)

				for i := 0; i < len(list.JSON200.Profiles)-1; i++ {
					a.LessOrEqual(
						strings.ToLower(list.JSON200.Profiles[i].Name),
						strings.ToLower(list.JSON200.Profiles[i+1].Name),
						"profiles should be sorted by name in ascending order",
					)
				}
			})

			t.Run("sort_by_name_descending", func(t *testing.T) {
				sortParam := "-name"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Sort: &sortParam,
					Q:    &prefix,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)

				for i := 0; i < len(list.JSON200.Profiles)-1; i++ {
					a.GreaterOrEqual(
						strings.ToLower(list.JSON200.Profiles[i].Name),
						strings.ToLower(list.JSON200.Profiles[i+1].Name),
						"profiles should be sorted by name in descending order",
					)
				}
			})

			t.Run("sort_by_created_at", func(t *testing.T) {
				sortParam := "created_at"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Sort: &sortParam,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)
			})

			t.Run("filter_by_join_date_after", func(t *testing.T) {
				joinedParam := "2020-01-01T00:00:00Z/"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Joined: &joinedParam,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)
			})

			t.Run("filter_by_join_date_before", func(t *testing.T) {
				joinedParam := "/2050-01-01T00:00:00Z"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Joined: &joinedParam,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)
			})

			t.Run("filter_by_join_date_range", func(t *testing.T) {
				joinedParam := "2020-01-01T00:00:00Z/2050-01-01T00:00:00Z"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Joined: &joinedParam,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)
			})

			t.Run("filter_by_join_date_no_slash", func(t *testing.T) {
				joinedParam := "2020-01-01T00:00:00Z"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Joined: &joinedParam,
				})
				tests.Ok(t, err, list)
				r.GreaterOrEqual(len(list.JSON200.Profiles), 3)
			})

			t.Run("empty_result_future_date", func(t *testing.T) {
				joinedParam := "2099-01-01T00:00:00Z"
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Joined: &joinedParam,
				})
				tests.Ok(t, err, list)
				a.Empty(list.JSON200.Profiles)
			})

			t.Run("invalid_time_range_format", func(t *testing.T) {
				joinedParam := "2020-01-01T00:00:00Z/2025-01-01T00:00:00Z/2050-01-01T00:00:00Z"
				tests.AssertRequest(cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Joined: &joinedParam,
				}))(t, 400)
			})

			_ = acc1
			_ = acc2
			_ = acc3
		}))
	}))
}

func TestProfileListRoleFilter(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		ar *account_querier.Querier,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, account.Account{
				Handle: "admin",
				Name:   "Admin User",
				Admin:  true,
			})
			adminSession := sh.WithSession(adminCtx)

			role1Name := "moderator-" + xid.New().String()
			role1, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
				Name:        role1Name,
				Colour:      "blue",
				Permissions: openapi.PermissionList{openapi.MANAGECATEGORIES},
			}, adminSession)
			tests.Ok(t, err, role1)

			role2Name := "editor-" + xid.New().String()
			role2, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
				Name:        role2Name,
				Colour:      "green",
				Permissions: openapi.PermissionList{openapi.SUBMITLIBRARYNODE},
			}, adminSession)
			tests.Ok(t, err, role2)

			user1 := newAccount(t, root, cl, ar, "user1")
			user2 := newAccount(t, root, cl, ar, "user2")
			user3 := newAccount(t, root, cl, ar, "user3")

			_, err = cl.AccountAddRoleWithResponse(adminCtx, user1.Handle, role1.JSON200.Id, adminSession)
			r.NoError(err)

			_, err = cl.AccountAddRoleWithResponse(adminCtx, user2.Handle, role1.JSON200.Id, adminSession)
			r.NoError(err)
			_, err = cl.AccountAddRoleWithResponse(adminCtx, user2.Handle, role2.JSON200.Id, adminSession)
			r.NoError(err)

			_, err = cl.AccountAddRoleWithResponse(adminCtx, user3.Handle, role2.JSON200.Id, adminSession)
			r.NoError(err)

			t.Run("filter_by_single_role", func(t *testing.T) {
				roles := []openapi.Identifier{openapi.Identifier(role1.JSON200.Id)}
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Roles: &roles,
				})
				tests.Ok(t, err, list)

				foundUser1 := false
				foundUser2 := false
				for _, profile := range list.JSON200.Profiles {
					if profile.Handle == user1.Handle {
						foundUser1 = true
					}
					if profile.Handle == user2.Handle {
						foundUser2 = true
					}
				}
				a.True(foundUser1, "user1 should be in results")
				a.True(foundUser2, "user2 should be in results")
			})

			t.Run("filter_by_multiple_roles_conjunction", func(t *testing.T) {
				roles := []openapi.Identifier{
					openapi.Identifier(role1.JSON200.Id),
					openapi.Identifier(role2.JSON200.Id),
				}
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					Roles: &roles,
				})
				tests.Ok(t, err, list)

				foundUser2 := false
				foundUser1 := false
				foundUser3 := false
				for _, profile := range list.JSON200.Profiles {
					if profile.Handle == user1.Handle {
						foundUser1 = true
					}
					if profile.Handle == user2.Handle {
						foundUser2 = true
					}
					if profile.Handle == user3.Handle {
						foundUser3 = true
					}
				}
				a.True(foundUser2, "user2 has both roles and should be in results")
				a.False(foundUser1, "user1 has only role1, not both")
				a.False(foundUser3, "user3 has only role2, not both")
			})

			_ = admin
		}))
	}))
}

func TestProfileListInvitedByFilter(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		ar *account_querier.Querier,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			inviter1Ctx, inviter1 := e2e.WithAccount(root, aw, account.Account{
				Handle: "inviter1",
				Name:   "Inviter One",
				Admin:  true,
			})
			inviter1Session := sh.WithSession(inviter1Ctx)

			inviter2Ctx, inviter2 := e2e.WithAccount(root, aw, account.Account{
				Handle: "inviter2",
				Name:   "Inviter Two",
				Admin:  true,
			})
			inviter2Session := sh.WithSession(inviter2Ctx)

			message := "Join!"
			invite1, err := cl.InvitationCreateWithResponse(root, openapi.InvitationInitialProps{
				Message: &message,
			}, inviter1Session)
			tests.Ok(t, err, invite1)

			invite2, err := cl.InvitationCreateWithResponse(root, openapi.InvitationInitialProps{
				Message: &message,
			}, inviter2Session)
			tests.Ok(t, err, invite2)

			invitee1Handle := "invitee1-" + xid.New().String()
			invitee1Resp, err := cl.AuthPasswordSignupWithResponse(root, &openapi.AuthPasswordSignupParams{
				InvitationId: &invite1.JSON200.Id,
			}, openapi.AuthPair{Identifier: invitee1Handle, Token: "password"})
			tests.Ok(t, err, invitee1Resp)

			invitee2Handle := "invitee2-" + xid.New().String()
			invitee2Resp, err := cl.AuthPasswordSignupWithResponse(root, &openapi.AuthPasswordSignupParams{
				InvitationId: &invite2.JSON200.Id,
			}, openapi.AuthPair{Identifier: invitee2Handle, Token: "password"})
			tests.Ok(t, err, invitee2Resp)

			t.Run("filter_by_single_inviter", func(t *testing.T) {
				inviters := []openapi.AccountHandle{openapi.AccountHandle(inviter1.Handle)}
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					InvitedBy: &inviters,
				})
				tests.Ok(t, err, list)

				foundInvitee1 := false
				foundInvitee2 := false
				for _, profile := range list.JSON200.Profiles {
					if profile.Handle == invitee1Handle {
						foundInvitee1 = true
					}
					if profile.Handle == invitee2Handle {
						foundInvitee2 = true
					}
				}
				a.True(foundInvitee1, "invitee1 should be in results")
				a.False(foundInvitee2, "invitee2 was not invited by inviter1")
			})

			t.Run("filter_by_multiple_inviters", func(t *testing.T) {
				inviters := []openapi.AccountHandle{
					openapi.AccountHandle(inviter1.Handle),
					openapi.AccountHandle(inviter2.Handle),
				}
				list, err := cl.ProfileListWithResponse(root, &openapi.ProfileListParams{
					InvitedBy: &inviters,
				})
				tests.Ok(t, err, list)

				foundInvitee1 := false
				foundInvitee2 := false
				for _, profile := range list.JSON200.Profiles {
					if profile.Handle == invitee1Handle {
						foundInvitee1 = true
					}
					if profile.Handle == invitee2Handle {
						foundInvitee2 = true
					}
				}
				a.True(foundInvitee1, "invitee1 should be in results")
				a.True(foundInvitee2, "invitee2 should be in results")
			})

			_ = r
		}))
	}))
}

func newAccount(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, ar *account_querier.Querier, handle string) account.Account {
	r := require.New(t)

	hand1 := handle + "-" + xid.New().String()

	response, err := cl.AuthPasswordSignupWithResponse(ctx, nil, openapi.AuthPair{
		Identifier: hand1,
		Token:      "password",
	})
	r.NoError(err)
	r.NotNil(response)
	r.Equal(http.StatusOK, response.StatusCode())

	acc, err := ar.GetByID(ctx, account.AccountID(utils.Must(xid.FromString(response.JSON200.Id))))
	r.NoError(err)
	r.NotNil(acc)

	return acc.Account
}

func findProfile(profiles []openapi.PublicProfile, handle string) (openapi.PublicProfile, bool) {
	return lo.Find(profiles, func(p openapi.PublicProfile) bool {
		return p.Handle == handle
	})
}

func findProfilesByNamePrefix(profiles []openapi.PublicProfile, prefix string) []openapi.PublicProfile {
	return lo.Filter(profiles, func(p openapi.PublicProfile, _ int) bool {
		return strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(prefix))
	})
}

func getProfileIndex(profiles []openapi.PublicProfile, handle string) int {
	for i, p := range profiles {
		if p.Handle == handle {
			return i
		}
	}
	return -1
}
